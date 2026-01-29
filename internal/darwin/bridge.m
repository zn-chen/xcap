#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import <CoreGraphics/CoreGraphics.h>
#include <stdlib.h>
#include <string.h>
#include "bridge.h"

// Helper: Copy NSString to C string buffer
static void copy_nsstring_to_buffer(NSString *str, char *buffer, size_t buffer_size) {
    if (str == nil) {
        buffer[0] = '\0';
        return;
    }
    const char *cstr = [str UTF8String];
    if (cstr) {
        strncpy(buffer, cstr, buffer_size - 1);
        buffer[buffer_size - 1] = '\0';
    } else {
        buffer[0] = '\0';
    }
}

#pragma mark - Monitor Functions

int owl_get_all_monitors(OwlMonitorInfo **monitors, int *count) {
    @autoreleasepool {
        // Get active display list
        uint32_t max_displays = 16;
        CGDirectDisplayID display_ids[16];
        uint32_t display_count = 0;

        CGError err = CGGetActiveDisplayList(max_displays, display_ids, &display_count);
        if (err != kCGErrorSuccess || display_count == 0) {
            *monitors = NULL;
            *count = 0;
            return OWL_ERR_NO_MONITORS;
        }

        // Allocate monitor info array
        OwlMonitorInfo *result = (OwlMonitorInfo *)calloc(display_count, sizeof(OwlMonitorInfo));
        if (result == NULL) {
            *monitors = NULL;
            *count = 0;
            return OWL_ERR_ALLOC_FAILED;
        }

        // Get NSScreen array for friendly names
        NSArray<NSScreen *> *screens = [NSScreen screens];

        for (uint32_t i = 0; i < display_count; i++) {
            CGDirectDisplayID display_id = display_ids[i];
            CGRect bounds = CGDisplayBounds(display_id);

            result[i].id = display_id;
            result[i].x = (int32_t)bounds.origin.x;
            result[i].y = (int32_t)bounds.origin.y;
            result[i].width = (uint32_t)bounds.size.width;
            result[i].height = (uint32_t)bounds.size.height;

            // Find matching NSScreen for friendly name
            NSString *name = nil;
            for (NSScreen *screen in screens) {
                NSDictionary *desc = [screen deviceDescription];
                NSNumber *screenNumber = desc[@"NSScreenNumber"];
                if (screenNumber && [screenNumber unsignedIntValue] == display_id) {
                    name = [screen localizedName];
                    break;
                }
            }

            if (name) {
                copy_nsstring_to_buffer(name, result[i].name, sizeof(result[i].name));
            } else {
                snprintf(result[i].name, sizeof(result[i].name), "Display %u", display_id);
            }
        }

        *monitors = result;
        *count = (int)display_count;
        return OWL_OK;
    }
}

void owl_free_monitors(OwlMonitorInfo *monitors) {
    if (monitors) {
        free(monitors);
    }
}

int owl_capture_monitor(uint32_t display_id, OwlCaptureResult *result) {
    @autoreleasepool {
        // Get display bounds
        CGRect bounds = CGDisplayBounds(display_id);

        // Capture the screen
        CGImageRef image = CGWindowListCreateImage(
            bounds,
            kCGWindowListOptionAll,
            kCGNullWindowID,
            kCGWindowImageDefault
        );

        if (image == NULL) {
            return OWL_ERR_CAPTURE_FAILED;
        }

        // Get image properties
        size_t width = CGImageGetWidth(image);
        size_t height = CGImageGetHeight(image);
        size_t bytes_per_row = CGImageGetBytesPerRow(image);

        // Get pixel data
        CGDataProviderRef provider = CGImageGetDataProvider(image);
        CFDataRef data = CGDataProviderCopyData(provider);

        if (data == NULL) {
            CGImageRelease(image);
            return OWL_ERR_CAPTURE_FAILED;
        }

        const uint8_t *src = CFDataGetBytePtr(data);
        size_t data_length = CFDataGetLength(data);

        // Allocate and copy data
        result->data = (uint8_t *)malloc(data_length);
        if (result->data == NULL) {
            CFRelease(data);
            CGImageRelease(image);
            return OWL_ERR_ALLOC_FAILED;
        }

        memcpy(result->data, src, data_length);
        result->width = (uint32_t)width;
        result->height = (uint32_t)height;
        result->bytes_per_row = (uint32_t)bytes_per_row;
        result->data_length = (uint32_t)data_length;

        CFRelease(data);
        CGImageRelease(image);

        return OWL_OK;
    }
}

#pragma mark - Window Functions

int owl_get_all_windows(OwlWindowInfo **windows, int *count) {
    @autoreleasepool {
        // Get window list
        CFArrayRef window_list = CGWindowListCopyWindowInfo(
            kCGWindowListOptionOnScreenOnly | kCGWindowListExcludeDesktopElements,
            kCGNullWindowID
        );

        if (window_list == NULL) {
            *windows = NULL;
            *count = 0;
            return OWL_ERR_NO_WINDOWS;
        }

        CFIndex window_count = CFArrayGetCount(window_list);
        if (window_count == 0) {
            CFRelease(window_list);
            *windows = NULL;
            *count = 0;
            return OWL_ERR_NO_WINDOWS;
        }

        // First pass: count valid windows
        int valid_count = 0;
        for (CFIndex i = 0; i < window_count; i++) {
            CFDictionaryRef window_info = CFArrayGetValueAtIndex(window_list, i);

            // Check sharing state
            CFNumberRef sharing_state_ref = CFDictionaryGetValue(window_info, kCGWindowSharingState);
            int sharing_state = 0;
            if (sharing_state_ref) {
                CFNumberGetValue(sharing_state_ref, kCFNumberIntType, &sharing_state);
            }
            if (sharing_state == 0) {
                continue;
            }

            // Filter StatusIndicator
            CFStringRef name_ref = CFDictionaryGetValue(window_info, kCGWindowName);
            CFStringRef owner_ref = CFDictionaryGetValue(window_info, kCGWindowOwnerName);

            if (name_ref && owner_ref) {
                if (CFStringCompare(name_ref, CFSTR("StatusIndicator"), 0) == kCFCompareEqualTo &&
                    CFStringCompare(owner_ref, CFSTR("Window Server"), 0) == kCFCompareEqualTo) {
                    continue;
                }
            }

            valid_count++;
        }

        if (valid_count == 0) {
            CFRelease(window_list);
            *windows = NULL;
            *count = 0;
            return OWL_ERR_NO_WINDOWS;
        }

        // Allocate result array
        OwlWindowInfo *result = (OwlWindowInfo *)calloc(valid_count, sizeof(OwlWindowInfo));
        if (result == NULL) {
            CFRelease(window_list);
            *windows = NULL;
            *count = 0;
            return OWL_ERR_ALLOC_FAILED;
        }

        // Second pass: fill in window info
        int result_index = 0;
        for (CFIndex i = 0; i < window_count && result_index < valid_count; i++) {
            CFDictionaryRef window_info = CFArrayGetValueAtIndex(window_list, i);

            // Check sharing state
            CFNumberRef sharing_state_ref = CFDictionaryGetValue(window_info, kCGWindowSharingState);
            int sharing_state = 0;
            if (sharing_state_ref) {
                CFNumberGetValue(sharing_state_ref, kCFNumberIntType, &sharing_state);
            }
            if (sharing_state == 0) {
                continue;
            }

            // Filter StatusIndicator
            CFStringRef name_ref = CFDictionaryGetValue(window_info, kCGWindowName);
            CFStringRef owner_ref = CFDictionaryGetValue(window_info, kCGWindowOwnerName);

            if (name_ref && owner_ref) {
                if (CFStringCompare(name_ref, CFSTR("StatusIndicator"), 0) == kCFCompareEqualTo &&
                    CFStringCompare(owner_ref, CFSTR("Window Server"), 0) == kCFCompareEqualTo) {
                    continue;
                }
            }

            // Get window ID
            CFNumberRef window_id_ref = CFDictionaryGetValue(window_info, kCGWindowNumber);
            uint32_t window_id = 0;
            if (window_id_ref) {
                CFNumberGetValue(window_id_ref, kCFNumberIntType, &window_id);
            }

            // Get PID
            CFNumberRef pid_ref = CFDictionaryGetValue(window_info, kCGWindowOwnerPID);
            uint32_t pid = 0;
            if (pid_ref) {
                CFNumberGetValue(pid_ref, kCFNumberIntType, &pid);
            }

            // Get bounds
            CFDictionaryRef bounds_ref = CFDictionaryGetValue(window_info, kCGWindowBounds);
            CGRect bounds = CGRectZero;
            if (bounds_ref) {
                CGRectMakeWithDictionaryRepresentation(bounds_ref, &bounds);
            }

            result[result_index].id = window_id;
            result[result_index].pid = pid;
            result[result_index].x = (int32_t)bounds.origin.x;
            result[result_index].y = (int32_t)bounds.origin.y;
            result[result_index].width = (uint32_t)bounds.size.width;
            result[result_index].height = (uint32_t)bounds.size.height;

            // Get app name
            if (owner_ref) {
                NSString *owner = (__bridge NSString *)owner_ref;
                copy_nsstring_to_buffer(owner, result[result_index].app_name, sizeof(result[result_index].app_name));
            }

            // Get window title
            if (name_ref) {
                NSString *name = (__bridge NSString *)name_ref;
                copy_nsstring_to_buffer(name, result[result_index].title, sizeof(result[result_index].title));
            }

            result_index++;
        }

        CFRelease(window_list);

        *windows = result;
        *count = result_index;
        return OWL_OK;
    }
}

void owl_free_windows(OwlWindowInfo *windows) {
    if (windows) {
        free(windows);
    }
}

int owl_capture_window(uint32_t window_id, OwlCaptureResult *result) {
    @autoreleasepool {
        // Get window bounds first
        CFArrayRef window_list = CGWindowListCopyWindowInfo(
            kCGWindowListOptionIncludingWindow,
            window_id
        );

        if (window_list == NULL || CFArrayGetCount(window_list) == 0) {
            if (window_list) CFRelease(window_list);
            return OWL_ERR_NOT_FOUND;
        }

        CFDictionaryRef window_info = CFArrayGetValueAtIndex(window_list, 0);
        CFDictionaryRef bounds_ref = CFDictionaryGetValue(window_info, kCGWindowBounds);
        CGRect bounds = CGRectZero;
        if (bounds_ref) {
            CGRectMakeWithDictionaryRepresentation(bounds_ref, &bounds);
        }

        CFRelease(window_list);

        // Capture the window
        CGImageRef image = CGWindowListCreateImage(
            bounds,
            kCGWindowListOptionIncludingWindow,
            window_id,
            kCGWindowImageDefault
        );

        if (image == NULL) {
            return OWL_ERR_CAPTURE_FAILED;
        }

        // Get image properties
        size_t width = CGImageGetWidth(image);
        size_t height = CGImageGetHeight(image);
        size_t bytes_per_row = CGImageGetBytesPerRow(image);

        // Get pixel data
        CGDataProviderRef provider = CGImageGetDataProvider(image);
        CFDataRef data = CGDataProviderCopyData(provider);

        if (data == NULL) {
            CGImageRelease(image);
            return OWL_ERR_CAPTURE_FAILED;
        }

        const uint8_t *src = CFDataGetBytePtr(data);
        size_t data_length = CFDataGetLength(data);

        // Allocate and copy data
        result->data = (uint8_t *)malloc(data_length);
        if (result->data == NULL) {
            CFRelease(data);
            CGImageRelease(image);
            return OWL_ERR_ALLOC_FAILED;
        }

        memcpy(result->data, src, data_length);
        result->width = (uint32_t)width;
        result->height = (uint32_t)height;
        result->bytes_per_row = (uint32_t)bytes_per_row;
        result->data_length = (uint32_t)data_length;

        CFRelease(data);
        CGImageRelease(image);

        return OWL_OK;
    }
}

#pragma mark - Cleanup

void owl_free_capture_result(OwlCaptureResult *result) {
    if (result && result->data) {
        free(result->data);
        result->data = NULL;
    }
}
