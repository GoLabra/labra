import { createDayjsCustomValue } from "@/lib/utils/dayjs-custom-value";
import dayjs from 'dayjs';
import { localeConfig } from "@/config/locale-config";
import { arrayDistinct } from "@/lib/utils/array";


export const dateTimeToString = (value: any): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.dateTime.valueFormat);
};

export const dateTimeToDisplayString = (value: any): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.dateTime.displayFormat);
};

export const dateTimeToStringPrecision = (value: any): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.dateTime.valueFormatPrecision);
};

export const stringToDateTime = (value?: string | null): dayjs.Dayjs | null => {
    if(!value){
        return null;
    }
    if(value == "now()"){
        return createDayjsCustomValue(value);
    }

    if(value.length == 19) {
        const datetime = dayjs(value, localeConfig.dateTime.valueFormat);
        return datetime;
    }

    const time = dayjs(value, localeConfig.dateTime.valueFormatPrecision);
    return time;
}



export const dateToString = (value: any): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.date.valueFormat);
};

export const dateToDisplayString = (value: any): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.date.inputFormat);
};

export const stringToDate = (value?: string | null): dayjs.Dayjs | null => {
    if(!value){
        return null;
    }
    if(value == "now()"){
        return createDayjsCustomValue(value);
    }
    const time = dayjs(value, localeConfig.date.valueFormat);
    return time;
}



export const timeToStringPrecision = (value: any | null): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.time.valueFormatPrecision);
};

export const timeToString = (value: any | null): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.time.valueFormat);
};

export const timeToDisplayString = (value: any | null): string | null => {
    if (value == null) {
        return null;
    }

    if(value.$custom){
        return value.$custom;
    }

    return dayjs(value).format(localeConfig.time.inputFormat);
};

export const stringToTime = (value?: string): dayjs.Dayjs | null => {
    if(!value){
        return null;
    }
    if(value == "now()"){
        return createDayjsCustomValue(value);
    }

    if(value.length == 8) {
        const time = dayjs(value, localeConfig.time.valueFormat);
        return time;
    }

    const time = dayjs(value, localeConfig.time.valueFormatPrecision);
    return time;
}










// export type ValueConverter = {
//     pipe: (value: any, tags?: Record<string, any>) => any;
//     compose: (value: any, tags?: Record<string, any>) => any;
// }

// export const EmptyStringToNull: ValueConverter = ({
//     pipe: (value: any) => {
//         if (value == '') {
//             return null;
//         }
//         return value;
//     },
//     compose: (value?: string) => value ?? ''
// })

// export const DateTimeToString: ValueConverter = ({
//     pipe: (value: any) => {
//         if (value == null) {
//             return null;
//         }

//         if(value.$custom){
//             return value.$custom;
//         }

//         return dayjs(value).format(localeConfig.dateTime.valueFormat);
//     },
//     compose: (value?: string, tags?: Record<string, any>) => {
//         if(!value){
//             return value;
//         }
//         if(value == "now()"){
//             return createDayjsCustomValue(value);
//         }
//         const time = dayjs(value, localeConfig.dateTime.valueFormat);
//         return time;
//     }
// })

// export const DateToString: ValueConverter = ({
//     pipe: (value: any): string | null => {
//         if (value == null) {
//             return null;
//         }

//         if(value.$custom){
//             return value.$custom;
//         }

//         return dayjs(value).format(localeConfig.date.valueFormat);
//     },
//     compose: (value?: string): dayjs.Dayjs | null => {
//         if(!value){
//             return null;
//         }
//         if(value == "now()"){
//             return createDayjsCustomValue(value);
//         }
//         const time = dayjs(value, localeConfig.date.valueFormat);
//         return time;
//     }
// })

// export const TimeToString: ValueConverter = ({
//     pipe: (value: any) => {
//         if (value == null) {
//             return null;
//         }

//         if(value.$custom){
//             return value.$custom;
//         }

//         return dayjs(value).format(localeConfig.time.valueFormat);
//     },
//     compose: (value?: string) => {
//         if(!value){
//             return value;
//         }
//         if(value == "now()"){
//             return createDayjsCustomValue(value);
//         }
//         const time = dayjs(value, localeConfig.time.valueFormat);
//         return time;
//     }
// })

// export const NumberToString: ValueConverter = ({
//     pipe: (value?: number) => {
//         if(value == null || value == undefined){
//             return value;
//         }
//         return value.toString();
//     },
//     compose: (value?: string): number|null => {
//         if(value == null){
//             return null;
//         }
//         return +value;
//     }
// })

// export const BooleanToString: ValueConverter = ({
//     pipe: (value?: boolean) => {
//         if(value == null || value == undefined){
//             return value;
//         }
//         return value.toString();
//     },
//     compose: (value?: string): boolean | null => {
//         if(value == null){
//             return null;
//         }
//         return value.toLowerCase() === 'true';
//     }
// })

// export const ArrayToJSON: ValueConverter = ({
//     pipe: (value?: Array<any>) => {
//         return JSON.stringify(value);
//     },
//     compose: (value?: string): Array<any> | null => {
//         if(!value){
//             return null;
//         }
//         return JSON.parse(value);
//     }
// })

// export const StringToJSON: ValueConverter = ({
//     pipe: (value?: string,): Map<string, any> | null => {
//         if(!value){
//             return null;
//         }
//         return JSON.parse(value);
//     },
//     compose: (value?: string): string | null => {
//         if(!value){
//             return null;
//         }
//         return JSON.stringify(value, null, 2);
//     }
// })

// export const MultilineToArray: ValueConverter  = ({
//     pipe: (value?: string, tags?: Record<string, any>): Array<string> | null => {
//         if(value == null || value == undefined){
//             return null;
//         }
//         return arrayDistinct(value.split('\n'));
//     },
//     compose: (value?: Array<string>, tags?: Record<string, any>): string | null => {
//         if(value == null){
//             return null;
//         }
//         return value.join('\n');
//     }
// })

// export const LookupToGraphQlFormat: ValueConverter  = ({
//     pipe: (value?: any): any => {
//         return value?.value;
//     },
//     compose: (value: any, tags?: Record<string, any>): any => {
//         if(!value){
//             return null;
//         }

//         if ('label' in value && 'value' in value && ('connect' in value.value || 'create' in value.value)) {
//             return value;
//         }
        
//         return {
//             label: tags?.displayField ? value[tags.displayField] : value.id,
//             value: {
//                     connect: {
//                         id: value.id
//                     }
//                 }
//         }
//     }
// })

// export const pipe = (converter: ValueConverter | ValueConverter[], tags?: Record<string, any>) => {
//     const converters = Array.isArray(converter) ? converter : [converter];

//     return (value: any) => {
//         return converters.reduce((value, converter) => {
//             return converter.pipe(value, tags);
//         }, value);
//     }
// }

// export const compose = (converter: ValueConverter | ValueConverter[], tags?: Record<string, any>) => {
//     const converters = Array.isArray(converter) ? converter : [converter];

//     return (value: any) => {
//         return converters.reduce((value, converter) => {
//             return converter.compose(value, tags);
//         }, value);
//     }
// }