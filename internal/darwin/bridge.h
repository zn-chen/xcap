#ifndef XCAP_BRIDGE_H
#define XCAP_BRIDGE_H

#include <stdint.h>

// Error codes
#define XCAP_OK 0
#define XCAP_ERR_NO_MONITORS 1
#define XCAP_ERR_NO_WINDOWS 2
#define XCAP_ERR_CAPTURE_FAILED 3
#define XCAP_ERR_ALLOC_FAILED 4
#define XCAP_ERR_NOT_FOUND 5

// Monitor information
typedef struct {
    uint32_t id;
    char name[256];
    int32_t x;
    int32_t y;
    uint32_t width;
    uint32_t height;
} XcapMonitorInfo;

// Window information
typedef struct {
    uint32_t id;
    uint32_t pid;
    char app_name[256];
    char title[256];
    int32_t x;
    int32_t y;
    uint32_t width;
    uint32_t height;
} XcapWindowInfo;

// Capture result
typedef struct {
    uint8_t *data;
    uint32_t width;
    uint32_t height;
    uint32_t bytes_per_row;
    uint32_t data_length;
} XcapCaptureResult;

// Monitor functions
int xcap_get_all_monitors(XcapMonitorInfo **monitors, int *count);
void xcap_free_monitors(XcapMonitorInfo *monitors);
int xcap_capture_monitor(uint32_t display_id, XcapCaptureResult *result);

// Window functions
int xcap_get_all_windows(XcapWindowInfo **windows, int *count);
void xcap_free_windows(XcapWindowInfo *windows);
int xcap_capture_window(uint32_t window_id, XcapCaptureResult *result);

// Capture cleanup
void xcap_free_capture_result(XcapCaptureResult *result);

#endif // XCAP_BRIDGE_H
