#include <windows.h>
#include <dwmapi.h>
#include <shellscalingapi.h>
#include <psapi.h>
#include <stdlib.h>
#include <string.h>
#include "bridge.h"

// Link required libraries
#pragma comment(lib, "user32.lib")
#pragma comment(lib, "gdi32.lib")
#pragma comment(lib, "dwmapi.lib")
#pragma comment(lib, "shcore.lib")

// Constants (only define if not already defined by Windows headers)
#ifndef DWMWA_CLOAKED
#define DWMWA_CLOAKED 14
#endif
#ifndef DWMWA_EXTENDED_FRAME_BOUNDS
#define DWMWA_EXTENDED_FRAME_BOUNDS 9
#endif
#ifndef PW_RENDERFULLCONTENT
#define PW_RENDERFULLCONTENT 2
#endif

// ============================================================================
// Monitor Functions
// ============================================================================

typedef struct {
    XcapMonitorInfo *monitors;
    int count;
    int capacity;
} EnumMonitorData;

static BOOL CALLBACK monitor_enum_callback(HMONITOR hMonitor, HDC hdcMonitor,
                                           LPRECT lprcMonitor, LPARAM dwData) {
    EnumMonitorData *data = (EnumMonitorData *)dwData;

    // Expand array if needed
    if (data->count >= data->capacity) {
        int new_capacity = data->capacity * 2;
        XcapMonitorInfo *new_monitors = (XcapMonitorInfo *)realloc(
            data->monitors, new_capacity * sizeof(XcapMonitorInfo));
        if (new_monitors == NULL) {
            return FALSE;
        }
        data->monitors = new_monitors;
        data->capacity = new_capacity;
    }

    MONITORINFOEXW mi;
    mi.cbSize = sizeof(MONITORINFOEXW);

    if (GetMonitorInfoW(hMonitor, (LPMONITORINFO)&mi)) {
        XcapMonitorInfo *info = &data->monitors[data->count];
        info->handle = (uintptr_t)hMonitor;
        memcpy(info->name, mi.szDevice, sizeof(info->name));
        info->x = mi.rcMonitor.left;
        info->y = mi.rcMonitor.top;
        info->width = (uint32_t)(mi.rcMonitor.right - mi.rcMonitor.left);
        info->height = (uint32_t)(mi.rcMonitor.bottom - mi.rcMonitor.top);
        info->is_primary = (mi.dwFlags & MONITORINFOF_PRIMARY) != 0;
        data->count++;
    }

    return TRUE;
}

int xcap_get_all_monitors(XcapMonitorInfo **monitors, int *count) {
    EnumMonitorData data;
    data.capacity = 8;
    data.count = 0;
    data.monitors = (XcapMonitorInfo *)calloc(data.capacity, sizeof(XcapMonitorInfo));

    if (data.monitors == NULL) {
        *monitors = NULL;
        *count = 0;
        return XCAP_ERR_ALLOC_FAILED;
    }

    EnumDisplayMonitors(NULL, NULL, monitor_enum_callback, (LPARAM)&data);

    if (data.count == 0) {
        free(data.monitors);
        *monitors = NULL;
        *count = 0;
        return XCAP_ERR_NO_MONITORS;
    }

    *monitors = data.monitors;
    *count = data.count;
    return XCAP_OK;
}

void xcap_free_monitors(XcapMonitorInfo *monitors) {
    if (monitors) {
        free(monitors);
    }
}

int xcap_get_monitor_dpi(uintptr_t handle, uint32_t *dpi_x, uint32_t *dpi_y) {
    HRESULT hr = GetDpiForMonitor((HMONITOR)handle, MDT_EFFECTIVE_DPI, dpi_x, dpi_y);
    if (SUCCEEDED(hr)) {
        return XCAP_OK;
    }
    *dpi_x = 96;
    *dpi_y = 96;
    return XCAP_OK;
}

int xcap_capture_monitor(uintptr_t handle, int32_t x, int32_t y,
                         uint32_t width, uint32_t height, XcapCaptureResult *result) {
    if (width == 0 || height == 0) {
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Get desktop DC
    HWND hwnd_desktop = GetDesktopWindow();
    HDC hdc_screen = GetDC(hwnd_desktop);
    if (hdc_screen == NULL) {
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Create compatible DC
    HDC hdc_mem = CreateCompatibleDC(hdc_screen);
    if (hdc_mem == NULL) {
        ReleaseDC(hwnd_desktop, hdc_screen);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Create compatible bitmap
    HBITMAP hbitmap = CreateCompatibleBitmap(hdc_screen, (int)width, (int)height);
    if (hbitmap == NULL) {
        DeleteDC(hdc_mem);
        ReleaseDC(hwnd_desktop, hdc_screen);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    HGDIOBJ old_bitmap = SelectObject(hdc_mem, hbitmap);

    // Copy screen content
    BOOL blt_result = BitBlt(hdc_mem, 0, 0, (int)width, (int)height,
                              hdc_screen, x, y, SRCCOPY);

    if (!blt_result) {
        SelectObject(hdc_mem, old_bitmap);
        DeleteObject(hbitmap);
        DeleteDC(hdc_mem);
        ReleaseDC(hwnd_desktop, hdc_screen);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Prepare bitmap info
    BITMAPINFO bi;
    memset(&bi, 0, sizeof(bi));
    bi.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
    bi.bmiHeader.biWidth = (LONG)width;
    bi.bmiHeader.biHeight = -(LONG)height;  // Top-down DIB
    bi.bmiHeader.biPlanes = 1;
    bi.bmiHeader.biBitCount = 32;
    bi.bmiHeader.biCompression = BI_RGB;

    // Allocate pixel buffer
    uint32_t data_size = width * height * 4;
    uint8_t *pixel_data = (uint8_t *)malloc(data_size);
    if (pixel_data == NULL) {
        SelectObject(hdc_mem, old_bitmap);
        DeleteObject(hbitmap);
        DeleteDC(hdc_mem);
        ReleaseDC(hwnd_desktop, hdc_screen);
        return XCAP_ERR_ALLOC_FAILED;
    }

    // Get bitmap data
    int ret = GetDIBits(hdc_mem, hbitmap, 0, height, pixel_data, &bi, DIB_RGB_COLORS);

    SelectObject(hdc_mem, old_bitmap);
    DeleteObject(hbitmap);
    DeleteDC(hdc_mem);
    ReleaseDC(hwnd_desktop, hdc_screen);

    if (ret == 0) {
        free(pixel_data);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    result->data = pixel_data;
    result->width = width;
    result->height = height;
    result->data_length = data_size;

    return XCAP_OK;
}

// ============================================================================
// Window Functions
// ============================================================================

typedef struct {
    XcapWindowInfo *windows;
    int count;
    int capacity;
    DWORD current_pid;
    bool exclude_current_process;
} EnumWindowData;

static bool is_window_cloaked(HWND hwnd) {
    int cloaked = 0;
    HRESULT hr = DwmGetWindowAttribute(hwnd, DWMWA_CLOAKED, &cloaked, sizeof(cloaked));
    return SUCCEEDED(hr) && cloaked != 0;
}

static void get_extended_frame_bounds(HWND hwnd, RECT *rect) {
    RECT ext_rect;
    HRESULT hr = DwmGetWindowAttribute(hwnd, DWMWA_EXTENDED_FRAME_BOUNDS,
                                       &ext_rect, sizeof(ext_rect));
    if (SUCCEEDED(hr)) {
        *rect = ext_rect;
    }
}

static void get_process_name(DWORD pid, WCHAR *name_buf, int buf_size) {
    name_buf[0] = L'\0';
    HANDLE hProcess = OpenProcess(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, FALSE, pid);
    if (hProcess != NULL) {
        GetModuleBaseNameW(hProcess, NULL, name_buf, buf_size);
        CloseHandle(hProcess);
    }
}

static BOOL CALLBACK window_enum_callback(HWND hwnd, LPARAM lParam) {
    EnumWindowData *data = (EnumWindowData *)lParam;

    // Check visibility
    if (!IsWindowVisible(hwnd)) {
        return TRUE;
    }

    // Check if cloaked
    if (is_window_cloaked(hwnd)) {
        return TRUE;
    }

    // Get process ID
    DWORD pid = 0;
    GetWindowThreadProcessId(hwnd, &pid);

    // Optionally exclude current process
    if (data->exclude_current_process && pid == data->current_pid) {
        return TRUE;
    }

    // Get window style
    LONG_PTR ex_style = GetWindowLongPtrW(hwnd, GWL_EXSTYLE);

    // Get class name
    WCHAR class_name[256];
    GetClassNameW(hwnd, class_name, 256);

    // Filter tool windows (but keep taskbar)
    if (ex_style & WS_EX_TOOLWINDOW) {
        if (wcscmp(class_name, L"Shell_TrayWnd") != 0 &&
            wcscmp(class_name, L"Shell_SecondaryTrayWnd") != 0) {
            return TRUE;
        }
    }

    // Filter system windows
    if (wcscmp(class_name, L"Progman") == 0 ||
        wcscmp(class_name, L"Button") == 0 ||
        wcscmp(class_name, L"Windows.UI.Core.CoreWindow") == 0) {
        return TRUE;
    }

    // Get window rect
    RECT rect;
    if (!GetWindowRect(hwnd, &rect)) {
        return TRUE;
    }

    // Try to get DWM extended frame bounds
    get_extended_frame_bounds(hwnd, &rect);

    int width = rect.right - rect.left;
    int height = rect.bottom - rect.top;

    // Filter empty rects
    if (width <= 0 || height <= 0) {
        return TRUE;
    }

    // Expand array if needed
    if (data->count >= data->capacity) {
        int new_capacity = data->capacity * 2;
        XcapWindowInfo *new_windows = (XcapWindowInfo *)realloc(
            data->windows, new_capacity * sizeof(XcapWindowInfo));
        if (new_windows == NULL) {
            return FALSE;
        }
        data->windows = new_windows;
        data->capacity = new_capacity;
    }

    XcapWindowInfo *info = &data->windows[data->count];
    memset(info, 0, sizeof(XcapWindowInfo));

    info->handle = (uintptr_t)hwnd;
    info->pid = pid;
    info->x = rect.left;
    info->y = rect.top;
    info->width = (uint32_t)width;
    info->height = (uint32_t)height;

    // Get window title
    int title_len = GetWindowTextLengthW(hwnd);
    if (title_len > 0 && title_len < 255) {
        GetWindowTextW(hwnd, (LPWSTR)info->title, 256);
    }

    // Get process name
    get_process_name(pid, (LPWSTR)info->app_name, 260);

    data->count++;
    return TRUE;
}

int xcap_get_all_windows(XcapWindowInfo **windows, int *count, bool exclude_current_process) {
    EnumWindowData data;
    data.capacity = 32;
    data.count = 0;
    data.current_pid = GetCurrentProcessId();
    data.exclude_current_process = exclude_current_process;
    data.windows = (XcapWindowInfo *)calloc(data.capacity, sizeof(XcapWindowInfo));

    if (data.windows == NULL) {
        *windows = NULL;
        *count = 0;
        return XCAP_ERR_ALLOC_FAILED;
    }

    EnumWindows(window_enum_callback, (LPARAM)&data);

    if (data.count == 0) {
        free(data.windows);
        *windows = NULL;
        *count = 0;
        return XCAP_ERR_NO_WINDOWS;
    }

    *windows = data.windows;
    *count = data.count;
    return XCAP_OK;
}

void xcap_free_windows(XcapWindowInfo *windows) {
    if (windows) {
        free(windows);
    }
}

bool xcap_is_window_minimized(uintptr_t handle) {
    return IsIconic((HWND)handle) != 0;
}

bool xcap_is_window_maximized(uintptr_t handle) {
    return IsZoomed((HWND)handle) != 0;
}

bool xcap_is_window_focused(uintptr_t handle) {
    return GetForegroundWindow() == (HWND)handle;
}

int xcap_capture_window(uintptr_t handle, XcapCaptureResult *result) {
    HWND hwnd = (HWND)handle;

    // Get window rect
    RECT rect;
    if (!GetWindowRect(hwnd, &rect)) {
        return XCAP_ERR_CAPTURE_FAILED;
    }

    int width = rect.right - rect.left;
    int height = rect.bottom - rect.top;

    if (width <= 0 || height <= 0) {
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Get window DC
    HDC hdc_window = GetWindowDC(hwnd);
    if (hdc_window == NULL) {
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Create compatible DC
    HDC hdc_mem = CreateCompatibleDC(hdc_window);
    if (hdc_mem == NULL) {
        ReleaseDC(hwnd, hdc_window);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Create compatible bitmap
    HBITMAP hbitmap = CreateCompatibleBitmap(hdc_window, width, height);
    if (hbitmap == NULL) {
        DeleteDC(hdc_mem);
        ReleaseDC(hwnd, hdc_window);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    HGDIOBJ old_bitmap = SelectObject(hdc_mem, hbitmap);

    // Capture using fallback strategy (same as xcap Rust)
    BOOL captured = FALSE;

    // 1. Windows 8+ use PrintWindow with PW_RENDERFULLCONTENT
    if (!captured) {
        OSVERSIONINFOW osvi;
        memset(&osvi, 0, sizeof(osvi));
        osvi.dwOSVersionInfoSize = sizeof(osvi);
        // Use RtlGetVersion for accurate version on Windows 8.1+
        typedef NTSTATUS(WINAPI *RtlGetVersionPtr)(PRTL_OSVERSIONINFOW);
        HMODULE ntdll = GetModuleHandleW(L"ntdll.dll");
        if (ntdll) {
            RtlGetVersionPtr rtl_get_version = (RtlGetVersionPtr)GetProcAddress(ntdll, "RtlGetVersion");
            if (rtl_get_version) {
                rtl_get_version((PRTL_OSVERSIONINFOW)&osvi);
            }
        }
        if (osvi.dwMajorVersion >= 8 ||
            (osvi.dwMajorVersion == 6 && osvi.dwMinorVersion >= 2)) {
            captured = PrintWindow(hwnd, hdc_mem, PW_RENDERFULLCONTENT);
        }
    }

    // 2. DWM composition enabled: use PrintWindow with default flag
    if (!captured) {
        BOOL dwm_enabled = FALSE;
        if (SUCCEEDED(DwmIsCompositionEnabled(&dwm_enabled)) && dwm_enabled) {
            captured = PrintWindow(hwnd, hdc_mem, 0);
        }
    }

    // 3. Try PrintWindow with flag 4
    if (!captured) {
        captured = PrintWindow(hwnd, hdc_mem, 4);
    }

    // 4. Fall back to BitBlt
    if (!captured) {
        captured = BitBlt(hdc_mem, 0, 0, width, height, hdc_window, 0, 0, SRCCOPY);
    }

    if (!captured) {
        SelectObject(hdc_mem, old_bitmap);
        DeleteObject(hbitmap);
        DeleteDC(hdc_mem);
        ReleaseDC(hwnd, hdc_window);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    // Prepare bitmap info
    BITMAPINFO bi;
    memset(&bi, 0, sizeof(bi));
    bi.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
    bi.bmiHeader.biWidth = width;
    bi.bmiHeader.biHeight = -height;  // Top-down DIB
    bi.bmiHeader.biPlanes = 1;
    bi.bmiHeader.biBitCount = 32;
    bi.bmiHeader.biCompression = BI_RGB;

    // Allocate pixel buffer
    uint32_t data_size = (uint32_t)(width * height * 4);
    uint8_t *pixel_data = (uint8_t *)malloc(data_size);
    if (pixel_data == NULL) {
        SelectObject(hdc_mem, old_bitmap);
        DeleteObject(hbitmap);
        DeleteDC(hdc_mem);
        ReleaseDC(hwnd, hdc_window);
        return XCAP_ERR_ALLOC_FAILED;
    }

    // Get bitmap data
    int ret = GetDIBits(hdc_mem, hbitmap, 0, (UINT)height, pixel_data, &bi, DIB_RGB_COLORS);

    SelectObject(hdc_mem, old_bitmap);
    DeleteObject(hbitmap);
    DeleteDC(hdc_mem);
    ReleaseDC(hwnd, hdc_window);

    if (ret == 0) {
        free(pixel_data);
        return XCAP_ERR_CAPTURE_FAILED;
    }

    result->data = pixel_data;
    result->width = (uint32_t)width;
    result->height = (uint32_t)height;
    result->data_length = data_size;

    return XCAP_OK;
}

// ============================================================================
// Cleanup
// ============================================================================

void xcap_free_capture_result(XcapCaptureResult *result) {
    if (result && result->data) {
        free(result->data);
        result->data = NULL;
    }
}

// ============================================================================
// Utility
// ============================================================================

uint32_t xcap_get_os_major_version(void) {
    typedef NTSTATUS(WINAPI *RtlGetVersionPtr)(PRTL_OSVERSIONINFOW);
    HMODULE ntdll = GetModuleHandleW(L"ntdll.dll");
    if (ntdll) {
        RtlGetVersionPtr rtl_get_version = (RtlGetVersionPtr)GetProcAddress(ntdll, "RtlGetVersion");
        if (rtl_get_version) {
            RTL_OSVERSIONINFOW osvi;
            memset(&osvi, 0, sizeof(osvi));
            osvi.dwOSVersionInfoSize = sizeof(osvi);
            rtl_get_version(&osvi);
            return osvi.dwMajorVersion;
        }
    }
    return 6;  // Default to Windows Vista/7
}
