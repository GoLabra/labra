"use client";
import dynamic from "next/dynamic";
import { InputBaseComponentProps } from "@mui/material";
import { forwardRef } from "react";
const MDEditor = dynamic(() => import("@uiw/react-md-editor"), { ssr: false });
export const MarkdownEditor = forwardRef<
  HTMLDivElement,
  InputBaseComponentProps
>(function MarkdownEditor(props, ref) {
  return (
    <div ref={ref} style={{ width: "100%" }} data-color-mode="light">
      <MDEditor
        value={(props.value as string) ?? ""}
        onChange={(val) =>
          props.onChange?.({
            target: { name: props.name, value: val ?? "" },
          } as any)
        }
        onBlur={props.onBlur as any}
      />
    </div>
  );
});
