import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React, { useEffect, useMemo } from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { TextLongFormField } from "@/core-features/dynamic-form/form-fields/TextLongField";
import { MultipleChoiceFormField } from "@/core-features/dynamic-form/form-fields/MultipleChoice";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { ChangedFullEntity } from "@/types/entity";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { optionsFromArray } from "@/lib/utils/form-util";
import { arrayDistinct } from "@/lib/utils/array";
import { useFormContext } from "react-hook-form";
import { getValuesFromOptions } from "@/core-features/dynamic-form/form-fields/TagsSelectField";

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
    acceptedValues: z.string().nonempty("Accepted Values is required")
                                         .transform(val => val?.split(/\r?\n/) ?? [])
                                         .transform((val) => arrayDistinct(val)),
    defaultValue: z.array(z.any()).optional().nullable()
        .transform((val) => {
            if (!val || (Array.isArray(val) && val.length === 0)) {
                return undefined;
            }
            return JSON.stringify(val);
        })
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
}

export const toDefaultValue = (value: any): any => {
    if (!value) {
        return value;
    }

    return {
        ...value,
        acceptedValues: value.acceptedValues ? value.acceptedValues.join('\n') : undefined,
        defaultValue: value.defaultValue ? JSON.parse(value.defaultValue) : undefined
    }
}

export const DesignerForm = (props: DesignerFormProps) => {

    const rawAcceptedValues = props.formMethods.watch('acceptedValues');
    const acceptedValues: string[] = useMemo(() => rawAcceptedValues?.split(/\r?\n/) ?? [], [rawAcceptedValues]);
    const options = useMemo(() => optionsFromArray(arrayDistinct(acceptedValues)) ?? [], [acceptedValues]);

    // reset selected value if option removed
    const formContext = useFormContext();
    const defaultValue = props.formMethods.watch('defaultValue');
    useEffect(() => {
        if (!defaultValue) {
            return;
        }

        let refValue = getValuesFromOptions(defaultValue, options);

        if (JSON.stringify(defaultValue) === JSON.stringify(refValue)) {
            return;
        }

        if (refValue?.length) {
            formContext.setValue('defaultValue', refValue);
        } else {
            formContext.setValue('defaultValue', null);
        }
    }, [defaultValue, formContext.setValue, rawAcceptedValues]);

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
        <BooleanFormField name="required" label="Required" />
        <TextLongFormField name="acceptedValues" label="Accepted Values" required />
        <MultipleChoiceFormField name="defaultValue" label="Default Value" options={options} />
    </Stack>
    )
}
