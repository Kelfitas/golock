#include "keyboard.h"
#include <stdio.h>
#include <xcb/xcb.h>
#include <xcb/xproto.h>
#include <xcb/xkb.h>
#include <xcb/xtest.h>
#include <xcb/xcb_keysyms.h>
#include <xkbcommon/xkbcommon.h>
#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon-x11.h>
#include <xcb/xcb_aux.h>
#include <xcb/randr.h>

void sendEvent(unsigned int eventCode, char *event, unsigned int keyCode) {
    int screens;
    xcb_connection_t *conn;
    if (xcb_connection_has_error((conn = xcb_connect(NULL, &screens))) > 0) {
        return;
    }

    // xcb_screen_t *screen;
    // screen = xcb_setup_roots_iterator(xcb_get_setup(conn)).data;

    printf("String is %s \n" , event);
    printf("keyCode %d \n" , keyCode);
    printf("eventCode %d \n" , eventCode);
    printf("KEY_PRESS %d \n" , XCB_KEY_PRESS);
    printf("KEY_RELEASE %d \n" , XCB_KEY_RELEASE);
    
    // xcb_send_event(conn, 1, screen->root, eventCode, (char *)event);
    // xcb_test_fake_input(conn, XCB_KEY_PRESS, keyCode, XCB_CURRENT_TIME, XCB_NONE, 0, 0, 0);
    xcb_test_fake_input(conn, eventCode, keyCode, XCB_CURRENT_TIME, XCB_NONE, 0, 0, 0);
    xcb_flush(conn);
    xcb_disconnect(conn);
}
