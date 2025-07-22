"use client";

import {
  FilledInput,
  FormControl,
  FormHelperText,
  InputLabel,
} from "@mui/material";
import { get, useFormContext } from "react-hook-form";
import useId from "@mui/material/utils/useId";
import { useLiteController } from "../lite-controller";
import { useFormDynamicContext } from "@/core-features/dynamic-form2/dynamic-form";
import { MarkdownEditor } from "@/shared/components/markdown-editor";

interface MarkdownFormComponentProps {
  name: string;
  label: string;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  errors?: any;

  value: any;
  onChange: (event: any) => void;
  onBlur: (event: any) => void;
}

export function MarkdownFormComponent(props: MarkdownFormComponentProps) {
  const {
    name,
    label,
    placeholder,
    disabled,
    errors,
    value,
    onChange,
    onBlur,
  } = props;

  const id = useId();
  return (
    <FormControl fullWidth variant="filled" error={!!errors}>
      <InputLabel htmlFor={id} id={`${id}-label`}>
        {label}
      </InputLabel>
      <FilledInput
        fullWidth
        inputComponent={MarkdownEditor as any}
        placeholder={placeholder}
        disabled={disabled}
        id={id}
        name={name}
        value={value ?? ""}
        onChange={onChange}
        sx={{ overflow: "inherit", padding: 0 }}
      />
      {errors && <FormHelperText>{errors}</FormHelperText>}
    </FormControl>
  );
}

interface FormFieldProps {
  name: string;
  placeholder?: string;
  label: string;
  disabled?: boolean;
  hide?: boolean;
  required?: boolean;
}
export function MarkdownFormField(props: FormFieldProps) {
  useFormDynamicContext(props.name, { disabled: props.disabled });
  const formContext = useFormContext();
  const formControllerHandler = useLiteController({
    name: props.name,
    control: formContext.control,
    disabled: props.disabled,
  });
  useFormDynamicContext(props.name, {
    disabled: formControllerHandler.disabled,
  });

  if (props.hide) {
    return null;
  }

  const errors = get(formContext.formState.errors, props.name);

  return (
    <MarkdownFormComponent
      label={props.label}
      placeholder={props.placeholder}
      required={props.required}
      errors={errors?.message}
      {...formControllerHandler}
    />
  );
}
