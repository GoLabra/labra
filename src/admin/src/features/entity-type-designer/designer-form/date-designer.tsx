import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { DateFormField } from "@/core-features/dynamic-form/form-fields/DateField";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { ChangedFullEntity } from "@/types/entity";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { dateToString, stringToDate } from "@/core-features/dynamic-form/value-convertor";

export const schema = z.object({
    caption: z.string().nonempty('Caption is required')
        .refine((val) => {
            const normalizedCaption = val.replace(/[\s_]/g, "").toLocaleLowerCase();
            return ENTITY_CHILDREN_SYSTEM_KEYWORDS.includes(normalizedCaption) == false;
        }, {
            message: 'System reserved keyword',
        })
        .refine((val) => {
            if (!val) {
                return true;
            }
            const firstLetter = val.charAt(0);
            return firstLetter.toLocaleLowerCase() != firstLetter.toLocaleUpperCase();
        }, {
            message: 'Caption must start with a letter',
        }),
    required: z.coerce.boolean().optional(),
    unique: z.coerce.boolean().optional(),
    min: z.literal("").transform(() => undefined).or(z.any().transform((val) => dateToString(val))).optional(),
    max: z.literal("").transform(() => undefined).or(z.any().transform((val) => dateToString(val))).optional(),
    defaultValue: z.literal("").transform(() => undefined).or(z.any().transform((val) => dateToString(val))).optional()
});

export const schemaEffect = (schema: z.ZodObject<any, any>, editId: string | undefined, entity: ChangedFullEntity): z.ZodEffects<any, any> => {
    return schema.refine((data) => {
        if (data.caption == undefined) {
            return false;
        }
        return checkChildCaptionExists(data.caption, editId, entity);
    }, {
        path: ['caption'],
        message: 'Caption already exists'
    })
        .superRefine((data, ctx) => {
            if (data.defaultValue == undefined) {
                return;
            }
            if (data.min !== undefined && stringToDate(data.defaultValue)!.isBefore(stringToDate(data.min))) {
                ctx.addIssue({
                    code: z.ZodIssueCode.custom,
                    path: ["defaultValue"],
                    message: "defaultValue must be >= Min Value",
                });
            }
            if (data.max !== undefined && stringToDate(data.defaultValue)!.isAfter(stringToDate(data.max))) {
                ctx.addIssue({
                    code: z.ZodIssueCode.custom,
                    path: ["defaultValue"],
                    message: "defaultValue must be <= Max Value",
                });
            }
        });
}

export const toDefaultValue = (value: any): any => {
    if(!value){
        return value;
    }
    
    return {
        ...value,
        min: value.min ? stringToDate(value.min) : value.min,
        max: value.max ? stringToDate(value.max) : value.max,
        defaultValue: value.defaultValue ? stringToDate(value.defaultValue) : value.defaultValue
    }
}

export const DesignerForm = (props: DesignerFormProps) => {

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
        <BooleanFormField name="required" label="Required" />

        <Stack direction="row" gap={1}>
            <DateFormField name="min" label="Min Value"
                max={props.formMethods.watch('max')} />
            <DateFormField name="max" label="Max Value"
                min={props.formMethods.watch('min')} />
        </Stack>

        <DateFormField name="defaultValue" label="Default Value"
            disabled={props.formMethods.watch('unique') === true}
            min={props.formMethods.watch('min')}
            max={props.formMethods.watch('max')} />
    </Stack>
    )
}
