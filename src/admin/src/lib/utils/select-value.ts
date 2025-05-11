export const selectValue = (key: string, obj: any):any => {
    if(!key){
        return null;
    }

    if(!obj){
        return null;
    }

    const keyArray = key.split('.');
    const result = selectValueFromObject(keyArray, obj);
    return result;
}

const selectValueFromObject = (keyArray: string[], obj: any):any => {
    const key = keyArray.shift();
    if(!key){
        return null;
    }

    const value = obj[key];
    if(!value){
        return null;
    }

    if(!keyArray.length){
        return value;
    }

    return selectValueFromObject(keyArray, value);
}