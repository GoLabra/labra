import { Lock, LockOpen, TextFields } from "@mui/icons-material";
import { Box, Button, Divider, InputBaseComponentProps, Stack, Typography } from "@mui/material";
import type { Editor, EditorOptions } from "@tiptap/core";
import { forwardRef, PropsWithChildren, useCallback, useEffect, useImperativeHandle, useRef, useState } from "react";
import {
    DebounceRender,
    LinkBubbleMenu,
    MenuButton,
    RichTextContent,
    RichTextEditor,
    RichTextEditorProvider,
    RichTextReadOnly,
    TableBubbleMenu,
    insertImages,
    type RichTextEditorRef,
} from "mui-tiptap";
import EditorMenuControls, { getInlineFiles } from "./editorMenuControls";
import useExtensions from "./useExtensions";
import { useEditor } from "@tiptap/react";

// const exampleContent = 'hello';

function fileListToImageFiles(fileList: FileList): File[] {
    // You may want to use a package like attr-accept
    // (https://www.npmjs.com/package/attr-accept) to restrict to certain file
    // types.
    return Array.from(fileList).filter((file) => {
        const mimeType = (file.type || "").toLowerCase();
        return mimeType.startsWith("image/");
    });
}


export const TiptapEditor = forwardRef((props: InputBaseComponentProps, ref) => {

    const extensions = useExtensions({
        placeholder: props.placeholder
    });

    const handleNewImageFiles = useCallback(
        async (files: File[], insertPosition?: number) => {
            if (!editor) {
                return;
            }

            insertImages({
                images: await getInlineFiles(files),
                editor: editor,
                position: insertPosition
            });
        }, []);

    // Allow for dropping images into the editor
    const handleDrop: NonNullable<EditorOptions["editorProps"]["handleDrop"]> =
        useCallback(
            (view, event, _slice, _moved) => {
                if (!(event instanceof DragEvent) || !event.dataTransfer) {
                    return false;
                }

                const imageFiles = fileListToImageFiles(event.dataTransfer.files);
                if (imageFiles.length > 0) {
                    const insertPosition = view.posAtCoords({
                        left: event.clientX,
                        top: event.clientY,
                    })?.pos;

                    handleNewImageFiles(imageFiles, insertPosition);

                    // Return true to treat the event as handled. We call preventDefault
                    // ourselves for good measure.
                    event.preventDefault();
                    return true;
                }

                return false;
            },
            [handleNewImageFiles],
        );

    // Allow for pasting images
    const handlePaste: NonNullable<EditorOptions["editorProps"]["handlePaste"]> =
        useCallback(
            (_view, event, _slice) => {
                if (!event.clipboardData) {
                    return false;
                }

                const pastedImageFiles = fileListToImageFiles(
                    event.clipboardData.files,
                );
                if (pastedImageFiles.length > 0) {
                    handleNewImageFiles(pastedImageFiles);
                    // Return true to mark the paste event as handled. This can for
                    // instance prevent redundant copies of the same image showing up,
                    // like if you right-click and copy an image from within the editor
                    // (in which case it will be added to the clipboard both as a file and
                    // as HTML, which Tiptap would otherwise separately parse.)
                    return true;
                }

                // We return false here to allow the standard paste-handler to run.
                return false;
            },
            [handleNewImageFiles],
        );

    const editor = useEditor({
        extensions: extensions,
        content: props.value,
        editable: true,
        onUpdate: ({ editor }) => {

            const htmlText = editor.getHTML();
            
            if(!htmlText){
                props.onChange?.({ target: { name: props.name, value: '' } } as any);
            }
            
            const rawText = editor.getText().replace(/<!--|-->|/g, '')
                                            .replace(/\n+/g, ' ');

            const content = `<!-- ${rawText} --> ${htmlText}`;
            props.onChange?.({ target: { name: props.name, value: content } } as any);
        },
        // optional:
        editorProps: {
            handleDrop: handleDrop,
            handlePaste: handlePaste
        },
		immediatelyRender: false
    });

    useImperativeHandle(ref, () => ({
        focus: () => editor?.chain().focus(),
    }))

    if (!editor) return null; // editor might be null on the first render

    return (
        <Box width={1} sx={{ overflow: 'visible' }}>
            <RichTextEditorProvider editor={editor}>

                <MenuBar>
                    <Box sx={{padding: '3px',}}>
                        <DebounceRender>
                            <EditorMenuControls />
                        </DebounceRender>    
                    </Box>
                    <Divider />
                </MenuBar>
                

                <Box sx={{ padding: 'calc(1.5 * var(--mui-spacing))' }}>
                    <RichTextContent />
                      
                    <LinkBubbleMenu />
                    <TableBubbleMenu />
                </Box>

            </RichTextEditorProvider>
        </Box>
    );
})


const MenuBar = (props: PropsWithChildren) => {
    return (<Box sx={{
        position: "sticky",
        top: -20,
        zIndex: 2,
        backgroundColor: "background.paper",
        borderRadius: 1
    }}>

        {props.children}

    </Box>)
}

TiptapEditor.displayName = 'TiptapEditor';