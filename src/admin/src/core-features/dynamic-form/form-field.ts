
export type Option<T = string | number | boolean | null, TT = never> = {
    label: string;
    tag?: TT;
    value: T;
};

// export type ObjectOption<T> = {
//     label: string;
//     value: T;
// };

// export type OptionWithSetting <T>= {
//     options: Array<PrimitiveOption<T>>,
//     valueProp: string;
// }

export type Options<T = string | number | boolean | null, TT = never> = Option<T, TT>[]; // | OptionWithSetting<Record<string, any>>;

export enum FormOpenMode {
    New = 1,
    Edit = 2,
    View = 3
}