import { Options } from "@/core-features/dynamic-form/form-field";

export const optionsFromArray = (value: Array<string>): Options | null => {
    if(!Array.isArray(value)){
        return null;
    }
    return value.map(item => ({ label: item, value: item }));
}