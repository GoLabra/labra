import { useCallback, useMemo } from "react"

const saveAs = typeof window === 'undefined' ? undefined : require('save-as').default

export const useExport = (entityName: string) => {

    const exportData = useCallback((data: any[]) => {
        let blob = new Blob([JSON.stringify(data)], { type: 'text/json;charset=utf-8' })
        saveAs(blob, `${entityName}.json`);
    }, [entityName]);

    const importFile = useCallback((data: any[]) => {
        
    }, []);

    return useMemo(() => ({
        exportData: exportData,
        importFile: importFile
    }), [exportData, importFile]);
}