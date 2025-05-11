import dayjs, { ConfigType, ManipulateType, OpUnitType, UnitType } from "dayjs";

var utc = require("dayjs/plugin/utc");
var timezone = require("dayjs/plugin/timezone");

dayjs.extend(utc);
dayjs.extend(timezone)

export const createDayjsCustomValue = (value: any) => {

    let _jsValue = dayjs();
    (_jsValue as any).$custom = value;

    return _jsValue;
};

export const isDayjsCustomValue = (dayjsValue: any): boolean => {
    return !!dayjsValue?.$custom;
}
