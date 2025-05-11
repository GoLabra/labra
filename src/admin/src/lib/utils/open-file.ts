const saveAsLib = typeof window === 'undefined' ? undefined : require('save-as').default

export const saveAs = (blob: Blob, fileName: string): void => {
    saveAsLib(blob, fileName);
}

export const saveJsonAs = (fileName: string, data: any): void => {
    let blob = new Blob([JSON.stringify(data)], { type: 'text/json;charset=utf-8' })
    saveAs(blob, fileName);
}


export const openOneJsonFile = (): Promise<any> => {

    return new Promise((resolve, reject) => {

      openOneFile(['.json']).then(file => {
            file.text().then(text => {
                try {
                    const data = JSON.parse(text);
                    resolve(data);
                } catch (e) {
                    reject('Invalid JSON format');
                }
            });
        }).catch(_ => reject('No file selected'));
    });
}

export const openOneFile = (acceptedExtensions: Array<string>): Promise<File> => {

    return new Promise((resolve, reject) => {

        var fileInput = document.createElement("input");
        fileInput.accept = acceptedExtensions.join(",");
        fileInput.type = "file";
        document.body.appendChild(fileInput);
        fileInput.onchange = () => {

            var file = fileInput.files?.[0];

            if (!file) {
                reject("No file selected");
            }

            resolve(file!);
        };

        fileInput.click();
        document.body.removeChild(fileInput);
    });
}