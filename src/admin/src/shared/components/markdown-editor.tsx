"use client";
import { Box, InputBaseComponentProps } from "@mui/material";
import MDEditor from "@uiw/react-md-editor";
import { forwardRef, useImperativeHandle } from "react";

export const MarkdownEditor = forwardRef(function MarkdownEditor(
  props: InputBaseComponentProps,
  ref,
) {
  const { name, value, onChange } = props;

  useImperativeHandle(ref, () => ({
    focus: () => {
      /* no-op */
    },
  }));

  return (
    <Box width={1} sx={{ overflow: "visible" }}>
      <MDEditor
        value={(value as string) ?? ""}
        onChange={(val) => {
          onChange?.({ target: { name, value: val || "" } } as any);
        }}
      />
    </Box>
  );
});

MarkdownEditor.displayName = "MarkdownEditor";
