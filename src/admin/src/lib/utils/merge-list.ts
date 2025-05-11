export const mergeListByProperty = (propertyName: string, originalList: Array<any>, changedList: Array<any>) => {

    var merged = originalList.map(item => {

        const changedEntity = changedList.find(i => i[propertyName] == item[propertyName]);

        if (changedEntity) {
            return {
                ...changedEntity,
            }
        };

        return item;
    });

    changedList.forEach(item => {
        if (!merged.find(i => i[propertyName] == item[propertyName])) {
            return;
        }
        merged.push(item);
    });

    return merged;
}