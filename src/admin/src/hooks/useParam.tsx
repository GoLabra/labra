import { useParams } from "next/navigation";

export const useParamSingleValue = (key: string, defaultValue?: string): string | undefined => {

    const params = useParams();
    const value = params[key] as string | undefined;

    return value ?? defaultValue;
}