#ifndef XCAP_WINDOWS_BRIDGE_H
#define XCAP_WINDOWS_BRIDGE_H

#include <stdint.h>
#include <stdbool.h>

// Error codes
#define XCAP_OK 0
#define XCAP_ERR_NO_MONITORS 1
#define XCAP_ERR_NO_WINDOWS 2
#define XCAP_ERR_CAPTURE_FAILED 3
#define XCAP_ERR_ALLOC_FAILED 4
#define XCAP_ERR_NOT_FOUND 5

// Monitor information (using Windows native types)
typedef struct {
    uintptr_t handle;        // HMONITOR
    uint16_t  name[32];      // Device name (UTF-16)
    int32_t   x;
    int32_t   y;
    uint32_t  width;
    uint32_t  height;
    bool      is_primary;
} XcapMonitorInfo;

// Window information
typedef struct {
    uintptr_t handle;        // HWND
    uint32_t  pid;
    uint16_t  app_name[260]; // Process name (UTF-16)
    uint16_t  title[256];    // Window title (UTF-16)
    int32_t   x;
    int32_t   y;
    uint32_t  width;
    uint32_t  height;
} XcapWindowInfo;

// Capture result (BGRA pixel data)
typedef struct {
    uint8_t  *data;
    uint32_t  width;
    uint32_t  height;
    uint32_t  data_length;
} XcapCaptureResult;

// Monitor functions
int xcap_get_all_monitors(XcapMonitorInfo **monitors, int *count);
void xcap_free_monitors(XcapMonitorInfo *monitors);
int xcap_capture_monitor(uintptr_t handle, int32_t x, int32_t y,
                         uint32_t width, uint32_t height, XcapCaptureResult *result);
int xcap_get_monitor_dpi(uintptr_t handle, uint32_t *dpi_x, uint32_t *dpi_y);

// Window functions
int xcap_get_all_windows(XcapWindowInfo **windows, int *count, bool exclude_current_process);
void xcap_free_windows(XcapWindowInfo *windows);
int xcap_capture_window(uintptr_t handle, XcapCaptureResult *result);

// Window state functions
bool xcap_is_window_minimized(uintptr_t handle);
bool xcap_is_window_maximized(uintptr_t handle);
bool xcap_is_window_focused(uintptr_t handle);

// Capture cleanup
void xcap_free_capture_result(XcapCaptureResult *result);

// Utility
uint32_t xcap_get_os_major_version(void);

#endif // XCAP_WINDOWS_BRIDGE_H
