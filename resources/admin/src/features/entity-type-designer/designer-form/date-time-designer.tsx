import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { DateTimeFormField } from "@/core-features/dynamic-form/form-fields/DateTimeField";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { ChangedFullEntity } from "@/types/entity";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { dateTimeToString, stringToDateTime } from "@/core-features/dynamic-form/value-convertor";
import dayjs from "dayjs";

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
    min: z.literal("").transform(() => undefined).or(z.any().transform((val) => dateTimeToString(val))).optional(),
    max: z.literal("").transform(() => undefined).or(z.any().transform((val) => dateTimeToString(val))).optional(),
    defaultValue: z.literal("").transform(() => undefined).or(z.any().transform((val) => dateTimeToString(val))).optional()
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
            
            if (data.min !== undefined && stringToDateTime(data.defaultValue)!.isBefore(stringToDateTime(data.min))) {
                ctx.addIssue({
                    code: z.ZodIssueCode.custom,
                    path: ["defaultValue"],
                    message: "defaultValue must be >= Min Value",
                });
            }
            if (data.max !== undefined && stringToDateTime(data.defaultValue)!.isAfter(stringToDateTime(data.max))) {
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
        min: value.min ? stringToDateTime(value.min) : value.min,
        max: value.max ? stringToDateTime(value.max) : value.max,
        defaultValue: value.defaultValue ? stringToDateTime(value.defaultValue) : value.defaultValue
    }
}

export const DesignerForm = (props: DesignerFormProps) => {
 
    const min = props.formMethods.watch('min');
    const max = props.formMethods.watch('max');

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
        <BooleanFormField name="required" label="Required" />

        <Stack direction="row" gap={1}>
            <DateTimeFormField name="min" label="Min Value" max={max} />
            <DateTimeFormField name="max" label="Max Value" min={min} />
        </Stack>

        <DateTimeFormField name="defaultValue" label="Default Value"
            disabled={props.formMethods.watch('unique') === true} min={min} max={max} />
    </Stack>
    )
}
