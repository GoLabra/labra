
import { ChildTypeDescriptor } from "@/features/entity-type-designer/designer-field-map";
import { User } from "@/lib/apollo/graphql";
import { minify } from "next/dist/build/swc";
import { createContext, PropsWithChildren, useCallback, useContext, useEffect, useMemo, useRef, useState } from "react";
import { FormProvider, FormState, useFormContext, UseFormReturn } from "react-hook-form";
import { Schema } from "yup";
import { z } from "zod";

interface formProps {
    children: any;
    onSubmit: any;
    methods: any;
    dynamicSchema?: FormDynamicContextType
}
export function Form(props: PropsWithChildren<formProps>) {
    return (
            <FormDynamicContext.Provider value={props.dynamicSchema ?? null}>
        <FormProvider {...props.methods}>
                <form onSubmit={props.onSubmit} noValidate autoComplete="off">
                    {props.children}
                </form>
        </FormProvider>
            </FormDynamicContext.Provider>
    );
}

export type FormDynamicContextType = {
    fieldsOptions: Record<string, SchemaFieldOptions>,
    pickedSchema: z.ZodObject<any, any> | z.ZodEffects<any, any>,
    disable: (name: string, disabled?: boolean) => void;
    min: (name: string, min?: number | Date) => void;
    max: (name: string, max?: number | Date) => void;
}

const FormDynamicContext = createContext<FormDynamicContextType | null>(null);


const triggerValidation = (name: string, disabled: boolean, formContext: UseFormReturn) => {
    if(disabled){
        formContext.clearErrors(name); 
        return;
    } 

    setTimeout(() => {
        const should = (formContext.formState.isSubmitted || formContext.formState.touchedFields[name]);
        if(should){
            formContext.trigger(name);
        }
    }, 0);
    
}

export const useFormDynamicContext = (name: string, options?: { min?: number, max?: number, disabled?: boolean }) => {
    const context = useContext(FormDynamicContext)
    const formContextRef = useRef<UseFormReturn>(null!);
    formContextRef.current = useFormContext();

    useEffect(() => {
        if(!context){
            return;
        }
        context.disable(name, options?.disabled);
        triggerValidation(name, options?.disabled ?? false, formContextRef.current!);
    }, [options?.disabled]);

    useEffect(() => {
        if(!context){
            return;
        }
        context.min(name, options?.min);
        triggerValidation(name, options?.disabled ?? false, formContextRef.current!);
    }, [options?.min]);

    useEffect(() => {
        if(!context){
            return;
        }
        context.max(name, options?.max);
        triggerValidation(name, options?.disabled ?? false, formContextRef.current!);
    }, [options?.max]);

    return context;
};

export function useSchemaOptions(schema: z.ZodObject<any, any>, schemaEffect?: (schema: z.ZodObject<any, any>) => z.ZodObject<any, any> | z.ZodEffects<any, any> ): FormDynamicContextType {  

    const [fieldsOptions, setFieldsOptions] = useState<Record<string, SchemaFieldOptions>>({});
    
    const disable = useCallback((name: string, disabled?: boolean) => {
        setFieldsOptions((prev: Record<string, SchemaFieldOptions>) => {
            if ((disabled ?? false) != (prev[name]?.disabled ?? false)) {
                return {
                    ...prev,
                    [name]: {
                        ...prev[name],
                        disabled: disabled
                    }
                }
            }
            return prev;
        });
        
    }, [setFieldsOptions]);

    const min = useCallback((name: string, min?: number | Date) => {
        setFieldsOptions((prev: Record<string, SchemaFieldOptions>) => {
            if ((min ?? null) != (prev[name]?.min ?? null)) {
                return {
                    ...prev,
                    [name]: {
                        ...prev[name],
                        min: min
                    }
                }
            }
            return prev;
        });
    }, [setFieldsOptions]);

    const max = useCallback((name: string, max?: number | Date) => {
        setFieldsOptions((prev: Record<string, SchemaFieldOptions>) => {
            if ((max ?? null) != (prev[name]?.max ?? null)) {
                return {
                    ...prev,
                    [name]: {
                        ...prev[name],
                        max: max
                    }
                }
            }
            return prev;
        });
    }, [setFieldsOptions]);

    const pickedSchema = useMemo(() => {
        const keys = zodKeys(schema);
        const pickedKeys = keys.filter(key => fieldsOptions[key]?.disabled !== true)
            .reduce((acc: any, key: string) => {
                acc[key] = true;
                return acc;
            }, {});

        const pickedSchema = schema.pick(pickedKeys);
        if(!schemaEffect){
            return pickedSchema;
        }
        return schemaEffect(pickedSchema) ?? pickedSchema;
        
    }, [schema, schemaEffect, fieldsOptions]);

    return useMemo(() => ({
        fieldsOptions,
        pickedSchema,
        disable,
        min,
        max
    }), [fieldsOptions, pickedSchema, disable, min, max ]);
}

export type SchemaFieldOptions = {
    disabled?: boolean;
    min?: number | Date;
    max?: number | Date;
}



// get zod object keys recursively
const zodKeys = <T extends z.ZodTypeAny>(schema: T): string[] => {
    // make sure schema is not null or undefined
    if (schema === null || schema === undefined) return [];
    // check if schema is nullable or optional
    if (schema instanceof z.ZodNullable || schema instanceof z.ZodOptional) return zodKeys(schema.unwrap());
    // check if schema is an array
    if (schema instanceof z.ZodArray) return zodKeys(schema.element);
    // check if schema is an object
    if (schema instanceof z.ZodObject) {
        // get key/value pairs from schema
        const entries = Object.entries(schema.shape);
        // loop through key/value pairs
        return entries.flatMap(([key, value]) => {
            // get nested keys
            const nested = value instanceof z.ZodType ? zodKeys(value).map(subKey => `${key}.${subKey}`) : [];
            // return nested keys
            return nested.length ? nested : key;
        });
    }
    // return empty array
    return [];
};
