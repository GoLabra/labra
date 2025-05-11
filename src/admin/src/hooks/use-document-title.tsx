import { useEffect, useState } from "react";

interface useDocumentTitleProps {
    title?: string;
}
export const useDocumentTitle = (props: useDocumentTitleProps) => {

    const [ title, setTitle] = useState<string|undefined>(props.title);

    useEffect(() => {
        const newTitle = title ?? '';
        document.title = newTitle;
    }, [title]);

    return {
        title,
        setTitle
    };
}