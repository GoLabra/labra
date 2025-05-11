import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { NumberFormField } from "@/core-features/dynamic-form/form-fields/NumberField";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { ChangedFullEntity } from "@/types/entity";

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
    required: z.coerce.boolean(),
    unique: z.coerce.boolean(),
    min: z.preprocess((value) => {
        if (value === null || value === undefined || value === "") {
            return undefined;
        }
        return value;
    }, z.coerce.number().nullish()),
    max: z.preprocess((value) => {
        if (value === null || value === undefined || value === "") {
            return undefined;
        }
        return value;
    }, z.coerce.number().nullish()),
    defaultValue: z.preprocess((value) => {
        if (value === null || value === undefined || value === "") {
            return undefined;
        }
        return value;
    }, z.coerce.number().nullish())
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
    }).superRefine((data, ctx) => {
        if (data.defaultValue == undefined) {
            return;
        }
        if (data.min !== undefined && data.defaultValue < data.min) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                path: ["defaultValue"],
                message: "defaultValue must be >= Min Value",
            });
        }
        if (data.max !== undefined && data.defaultValue > data.max) {
            ctx.addIssue({
                code: z.ZodIssueCode.custom,
                path: ["defaultValue"],
                message: "defaultValue must be <= Max Value",
            });
        }
    });
}

export const toDefaultValue = (value: any): any => {
    return value;
}

export const DesignerForm = (props: DesignerFormProps) => {

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
        <Stack direction="row" gap={1}>
            <BooleanFormField name="required" label="Required" />
            <BooleanFormField name="unique" label="Unique"
                disabled={!!props.formMethods.watch('defaultValue')} />
        </Stack>

        <Stack direction="row" gap={1}>
            <NumberFormField name="min" label="Min Value"
                max={props.formMethods.watch('max')} />
            <NumberFormField name="max" label="Max Value"
                min={props.formMethods.watch('min')} />
        </Stack>

        <NumberFormField name="defaultValue" label="Default Value"
            disabled={props.formMethods.watch('unique') === true}
            min={props.formMethods.watch('min')}
            max={props.formMethods.watch('max')} />
    </Stack>
    )
}
