#ifndef OWL_BRIDGE_H
#define OWL_BRIDGE_H

#include <stdint.h>

// Error codes
#define OWL_OK 0
#define OWL_ERR_NO_MONITORS 1
#define OWL_ERR_NO_WINDOWS 2
#define OWL_ERR_CAPTURE_FAILED 3
#define OWL_ERR_ALLOC_FAILED 4
#define OWL_ERR_NOT_FOUND 5

// Monitor information
typedef struct {
    uint32_t id;
    char name[256];
    int32_t x;
    int32_t y;
    uint32_t width;
    uint32_t height;
} OwlMonitorInfo;

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
} OwlWindowInfo;

// Capture result
typedef struct {
    uint8_t *data;
    uint32_t width;
    uint32_t height;
    uint32_t bytes_per_row;
    uint32_t data_length;
} OwlCaptureResult;

// Monitor functions
int owl_get_all_monitors(OwlMonitorInfo **monitors, int *count);
void owl_free_monitors(OwlMonitorInfo *monitors);
int owl_capture_monitor(uint32_t display_id, OwlCaptureResult *result);

// Window functions
int owl_get_all_windows(OwlWindowInfo **windows, int *count);
void owl_free_windows(OwlWindowInfo *windows);
int owl_capture_window(uint32_t window_id, OwlCaptureResult *result);

// Capture cleanup
void owl_free_capture_result(OwlCaptureResult *result);

#endif // OWL_BRIDGE_H
