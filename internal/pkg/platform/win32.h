#pragma once

#ifndef WIN32_LEAN_AND_MEAN
#define WIN32_LEAN_AND_MEAN 1
#endif
#include <windows.h>

LRESULT WindowProc(HWND window, UINT message, WPARAM wParam, LPARAM lParam);
LRESULT CALLBACK BridgeProc(HWND window, UINT message, WPARAM wParam,
                            LPARAM lParam) {
  return WindowProc(window, message, wParam, lParam);
}
