#ifndef CLIPBOARD_H
#define CLIPBOARD_H
#include <stdbool.h>

// 剪贴板内容获取函数
const char* getPlainText(void);
const char* getHtmlText(void);
const char* getImageBase64(void);


// 热键相关函数
//void RegisterHotKey(void);
//void RegisterHotKeyDynamic(unsigned int keyCode, unsigned int modifiers);
void RegisterHotKeyDynamic(unsigned int keyCode, unsigned int modifiers);
void InstallHotKeyHandler(void);
void NSAppActivateIgnoringOtherApps();  // ✅ 只声明，不要实现！
void UnregisterHotKey(unsigned int keyCode, unsigned int modifiers);

// 剪贴板内容写入
bool setClipboardContent(const char* content, const char* type);
bool setPlainText(const char* content);
bool setHtmlText(const char* content);
bool setImageBase64(const char* base64);


void simulateCmdC(void);

void HideDockIcon(void);
bool isAccessibilityEnabled(void); // 检查辅助功能权限
bool requestAccessibilityPermission(void); // 触发辅助功能权限申请


const char* VisionOCR(const char* path);//ocr
const char *VisionOCRFromMemory(const void *data, int size);
#endif