#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>
#import <stdlib.h>
#import <AppKit/AppKit.h>
#import <ApplicationServices/ApplicationServices.h>
#import <Vision/Vision.h>
#import <Foundation/Foundation.h>
#import "system.h"

// ============ 剪贴板内容 ============

//读取剪贴板内容
const char* getClipboardContent(const char* type) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSString *nsType = [NSString stringWithUTF8String:type];
    NSString *content = [pasteboard stringForType:nsType];
    if (!content) return NULL;
    return strdup([content UTF8String]);
}

const char* getPlainText() {
    return getClipboardContent("public.utf8-plain-text");
}

const char* getHtmlText() {
    return getClipboardContent("public.html");
}

const char* getImageBase64() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSArray *classes = [NSArray arrayWithObject:[NSImage class]];
    NSDictionary *options = [NSDictionary dictionary];
    NSArray *items = [pasteboard readObjectsForClasses:classes options:options];

    if ([items count] > 0) {
        NSImage *image = [items objectAtIndex:0];
        NSBitmapImageRep *bitmap = [[NSBitmapImageRep alloc] initWithData:[image TIFFRepresentation]];
        NSData *pngData = [bitmap representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
        NSString *base64String = [pngData base64EncodedStringWithOptions:0];
        return strdup([base64String UTF8String]);
    }
    return NULL;
}

// ============ 热键相关 ============
// 声明 Go 回调函数
extern void hotKeyCallback(unsigned int keyCode, unsigned int modifiers);
// 定义热键最大数量
#define MAX_HOTKEYS 10
// 保存注册的 HotKeyRef
static EventHotKeyRef hotKeyRefs[MAX_HOTKEYS];
static EventHotKeyID hotKeyIDs[MAX_HOTKEYS];
static int hotKeyCount = 0;
static int handlerInstalled = 0;

OSStatus MyHotKeyHandler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
    EventHotKeyID hkID;
    GetEventParameter(theEvent, kEventParamDirectObject, typeEventHotKeyID, NULL, sizeof(hkID), NULL, &hkID);

    unsigned int keyCode = hkID.id >> 16;
    unsigned int modifiers = hkID.id & 0xFFFF;

    hotKeyCallback(keyCode, modifiers); // ✅ 正确传值给 Go
    return noErr;
}

void InstallHotKeyHandler() {
    if (handlerInstalled) return;

    EventTypeSpec eventType;
    eventType.eventClass = kEventClassKeyboard;
    eventType.eventKind = kEventHotKeyPressed;

    InstallApplicationEventHandler(&MyHotKeyHandler, 1, &eventType, NULL, NULL);
    handlerInstalled = 1;
}

void RegisterHotKeyDynamic(unsigned int keyCode, unsigned int modifiers) {
    InstallHotKeyHandler();

    if (hotKeyCount >= MAX_HOTKEYS) return;

    EventHotKeyID hkID;
    hkID.signature = 'htk1';
    hkID.id = (keyCode << 16) | (modifiers & 0xFFFF); // ✅ 编码 keyCode 和 modifiers

    RegisterEventHotKey(keyCode, modifiers, hkID, GetApplicationEventTarget(), 0, &hotKeyRefs[hotKeyCount]);
    // ✅ 保存 ID
    hotKeyIDs[hotKeyCount] = hkID;
    hotKeyCount++;
}
void NSAppActivateIgnoringOtherApps() {
    [[NSRunningApplication currentApplication] activateWithOptions:NSApplicationActivateIgnoringOtherApps];
}

void UnregisterHotKey(unsigned int keyCode, unsigned int modifiers) {
    unsigned int targetID = (keyCode << 16) | (modifiers & 0xFFFF);

    for (int i = 0; i < hotKeyCount; i++) {
        unsigned int existingID = hotKeyIDs[i].id;
        if (existingID == targetID) {
            // 注销
            UnregisterEventHotKey(hotKeyRefs[i]);

            // 清理数组（用后面的项覆盖）
            for (int j = i; j < hotKeyCount - 1; j++) {
                hotKeyRefs[j] = hotKeyRefs[j + 1];
                hotKeyIDs[j] = hotKeyIDs[j + 1];
            }
            hotKeyCount--;
            return;
        }
    }
}

// EventHotKeyRef gHotKeyRef;
// EventHotKeyID gHotKeyID;

// extern void hotKeyCallback(void);

// OSStatus MyHotKeyHandler(EventHandlerCallRef nextHandler, EventRef theEvent, void *userData) {
//     hotKeyCallback();
//     return noErr;
// }

// void RegisterHotKey() {
//     gHotKeyID.signature = 'htk1';
//     gHotKeyID.id = 1;

//     EventTypeSpec eventType;
//     eventType.eventClass = kEventClassKeyboard;
//     eventType.eventKind = kEventHotKeyPressed;

//     InstallApplicationEventHandler(&MyHotKeyHandler, 1, &eventType, NULL, NULL);
//     RegisterEventHotKey(kVK_ANSI_O, cmdKey | shiftKey, gHotKeyID, GetApplicationEventTarget(), 0, &gHotKeyRef);
// }

// // 实现动态注册热键的函数
// void RegisterHotKeyDynamic(unsigned int keyCode, unsigned int modifiers) {
//     gHotKeyID.signature = 'htk1';
//     gHotKeyID.id = 1;

//     EventTypeSpec eventType;
//     eventType.eventClass = kEventClassKeyboard;
//     eventType.eventKind = kEventHotKeyPressed;

//     InstallApplicationEventHandler(&MyHotKeyHandler, 1, &eventType, NULL, NULL);
//     RegisterEventHotKey(keyCode, modifiers, gHotKeyID, GetApplicationEventTarget(), 0, &gHotKeyRef);
// }

// ============ 写入剪贴板 ============
bool setClipboardContent(const char* content, const char* type) {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    [pasteboard clearContents];

    NSString *nsContent = [NSString stringWithUTF8String:content];
    NSString *nsType = [NSString stringWithUTF8String:type];
    NSArray *types = [NSArray arrayWithObject:nsType];
    [pasteboard declareTypes:types owner:nil];
    BOOL success = [pasteboard setString:nsContent forType:nsType];
    return success;
}

bool setPlainText(const char* content) {
    return setClipboardContent(content, "public.utf8-plain-text");
}

// bool setHtmlText(const char* content) {
//     return setClipboardContent(content, "public.html");
// }
bool setImageBase64(const char* base64) {
     if (!base64) return false;

    @autoreleasepool {
        NSString *base64Str = [NSString stringWithUTF8String:base64];
        NSData *imageData = [[NSData alloc] initWithBase64EncodedString:base64Str options:0];
        if (!imageData) return false;

        NSImage *image = [[NSImage alloc] initWithData:imageData];
        if (!image) return false;

        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];
        [pasteboard declareTypes:@[NSPasteboardTypePNG] owner:nil];

        return [pasteboard setData:imageData forType:NSPasteboardTypePNG];
    }
}


bool setHtmlText(const char* content) {
    if (!content) return false;
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSString *html = [NSString stringWithUTF8String:content];

    // 提取纯文本（你可以简单去掉 HTML 标签，或者传入另一个值）
    NSString *plain = [[html stringByReplacingOccurrencesOfString:@"<[^>]+>" withString:@"" options:NSRegularExpressionSearch range:NSMakeRange(0, html.length)] stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceAndNewlineCharacterSet]];

    [pasteboard clearContents];
    [pasteboard declareTypes:@[@"public.html", NSPasteboardTypeString] owner:nil];

    [pasteboard setString:html forType:@"public.html"];
    [pasteboard setString:plain forType:NSPasteboardTypeString];
    return YES;
}


//
void simulateCmdC(void) {
    CGEventRef cmdDown = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)0x37, true); // Cmd down
    CGEventRef cDown = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)8, true);      // C down
    CGEventSetFlags(cDown, kCGEventFlagMaskCommand);

    CGEventRef cUp = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)8, false);       // C up
    CGEventSetFlags(cUp, kCGEventFlagMaskCommand);

    CGEventRef cmdUp = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)0x37, false);  // Cmd up

    CGEventPost(kCGHIDEventTap, cmdDown);
    CGEventPost(kCGHIDEventTap, cDown);
    CGEventPost(kCGHIDEventTap, cUp);
    CGEventPost(kCGHIDEventTap, cmdUp);

    CFRelease(cmdDown);
    CFRelease(cDown);
    CFRelease(cUp);
    CFRelease(cmdUp);
}

void HideDockIcon(void) {
    NSApplication *app = [NSApplication sharedApplication];
    [app setActivationPolicy:NSApplicationActivationPolicyProhibited];
}

//
bool isAccessibilityEnabled() {
    return AXIsProcessTrusted();
}

// bool triggerAX() {
//     // 会触发权限申请弹窗
//     requestAccessibilityPermission();
//     return AXIsProcessTrusted();
// }

bool requestAccessibilityPermission() {
    const void *keys[] = { kAXTrustedCheckOptionPrompt };
    const void *values[] = { kCFBooleanTrue };
    CFDictionaryRef options = CFDictionaryCreate(kCFAllocatorDefault, keys, values, 1,
                                                 &kCFCopyStringDictionaryKeyCallBacks,
                                                 &kCFTypeDictionaryValueCallBacks);
    bool result = AXIsProcessTrustedWithOptions(options);
    CFRelease(options);
    return result;
}

//ocr
const char* VisionOCR(const char* path) {
    @autoreleasepool {
        NSString* imagePath = [NSString stringWithUTF8String:path];
        NSData* imageData = [NSData dataWithContentsOfFile:imagePath];
        NSImage* image = [[NSImage alloc] initWithData:imageData];
        if (!image) return "Image load failed";

        CGImageRef cgRef = [image CGImageForProposedRect:NULL context:nil hints:nil];
        if (!cgRef) return "CGImage conversion failed";

        VNRecognizeTextRequest *request = [[VNRecognizeTextRequest alloc] init];
        request.recognitionLevel = VNRequestTextRecognitionLevelAccurate;

        VNImageRequestHandler *handler = [[VNImageRequestHandler alloc] initWithCGImage:cgRef options:@{}];
        NSError *error = nil;
        [handler performRequests:@[request] error:&error];

        if (error) return [[error.localizedDescription UTF8String] copy];

        NSMutableString *result = [NSMutableString string];
        for (VNRecognizedTextObservation *obs in request.results) {
            VNRecognizedText *topCandidate = [[obs topCandidates:1] firstObject];
            if (topCandidate) {
                [result appendString:topCandidate.string];
                [result appendString:@"\n"];
            }
        }

        return strdup([result UTF8String]);
    }
}

const char *VisionOCRFromMemory(const void *data, int size) {
    @autoreleasepool {
        NSData *imageData = [NSData dataWithBytes:data length:size];
        NSImage *image = [[NSImage alloc] initWithData:imageData];
        CGImageRef cgRef = [image CGImageForProposedRect:nil context:nil hints:nil];
        VNRecognizeTextRequest *req = [[VNRecognizeTextRequest alloc] init];
        req.recognitionLevel = VNRequestTextRecognitionLevelAccurate;
        req.recognitionLanguages = @[@"zh-Hans", @"en-US"];
        VNImageRequestHandler *handler = [[VNImageRequestHandler alloc] initWithCGImage:cgRef options:@{}];
        NSError *err = nil;
        [handler performRequests:@[req] error:&err];
        if (err) return strdup([[err localizedDescription] UTF8String]);
        NSMutableString *out = [NSMutableString string];
        for (VNRecognizedTextObservation *obs in req.results) {
            VNRecognizedText *cand = [[obs topCandidates:1] firstObject];
            if (cand) {
                [out appendString:cand.string];
                [out appendString:@"\n"];
            }
        }
        return strdup([out UTF8String]);
    }
}