#pragma once

#include <windows.h>

#include <stdint.h>

#define internal static
#define global static

typedef int8_t i8;
typedef int16_t i16;
typedef int32_t i32;
typedef int64_t i64;

typedef uint8_t u8;
typedef uint16_t u16;
typedef uint32_t u32;
typedef uint64_t u64;

typedef float f32;
typedef double f64;

typedef i32 b32;

typedef u8 byte;

global i64 globalPerfCountFrequency = 0;

LRESULT WindowProc(HWND window, UINT message, WPARAM wParam, LPARAM lParam);
LRESULT CALLBACK BridgeProc(HWND window, UINT message, WPARAM wParam,
                            LPARAM lParam) {
  return WindowProc(window, message, wParam, lParam);
}

static inline LARGE_INTEGER Win32GetClockValue() {
  LARGE_INTEGER counter;
  QueryPerformanceCounter(&counter);
  return counter;
}

static inline f64 Win32GetSecondsElapsed(LARGE_INTEGER start,
                                         LARGE_INTEGER end) {
  return (f64)(end.QuadPart - start.QuadPart) / (f64)globalPerfCountFrequency;
}

static inline void Win32SetGlobalPerfFrequency() {
  LARGE_INTEGER perfCountFrequencyResult;
  QueryPerformanceFrequency(&perfCountFrequencyResult);
  globalPerfCountFrequency = perfCountFrequencyResult.QuadPart;
}

static inline b32 Win32SetSleepGranular() {
  UINT desiredSchedulerMs = 1;
  return (timeBeginPeriod(desiredSchedulerMs) == TIMERR_NOERROR);
}
