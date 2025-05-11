import { BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { Stack } from "@mui/material";
import React, { useEffect, useMemo } from "react";
import { z } from "zod";
import { DesignerFormProps } from "../designer-field-map";
import { BooleanSelectFormField } from "@/core-features/dynamic-form/form-fields/BooleanSelectField";
import { SingleChoiceFormField } from "@/core-features/dynamic-form/form-fields/SingleChoice";
import { TextLongFormField } from "@/core-features/dynamic-form/form-fields/TextLongField";
import { ENTITY_CHILDREN_SYSTEM_KEYWORDS } from "@/config/CONST";
import { ChangedFullEntity } from "@/types/entity";
import { checkChildCaptionExists } from "../use-designer-system-validation";
import { optionsFromArray } from "@/lib/utils/form-util";
import { useFormContext } from "react-hook-form";
import { getValuesFromOptions } from "@/core-features/dynamic-form/form-fields/TagsSelectField";
import { getValueFromOptions } from "@/core-features/dynamic-form/form-fields/SelectField";
import { arrayDistinct } from "@/lib/utils/array";

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
    defaultValue: z.preprocess((value) => value ? value : null, z.string().nullable()),
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
        acceptedValues: value.acceptedValues ? value.acceptedValues.join('\n') : value.acceptedValues,
    }
}

export const DesignerForm = (props: DesignerFormProps) => {

    const acceptedValues: string[] = props.formMethods.watch('acceptedValues')?.split(/\r?\n/) ?? [];
    const options = useMemo(() => optionsFromArray(arrayDistinct(acceptedValues)) ?? [], [acceptedValues]); 
    
    // reset selected value if option removed
    const formContext = useFormContext();
    const defaultValue = props.formMethods.watch('defaultValue');    
    useEffect(() => {
        if(!defaultValue){
            return;
        }

        let refValue = getValueFromOptions(defaultValue, options);
        if(refValue != null){
            return;
        }

        formContext.setValue('defaultValue', null);
    }, [defaultValue, formContext.setValue, options]);

    return (<Stack gap={1.5}>
        <TextShortFormField name="caption" label="Caption" required />
        <BooleanFormField name="required" label="Required" />
        <TextLongFormField name="acceptedValues" label="Accepted Values" required />
        <SingleChoiceFormField name="defaultValue" label="Default Value" options={options} />
    </Stack>
    )
}
