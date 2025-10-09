import i18next from 'i18next';
import { allLangs } from './all-langs';

export const fallbackLng = 'en';

export const localeConfig = {

    dateTime: {
        valueFormatPrecision: "YYYY-MM-DDTHH:mm:ss.SSS",
        valueFormat: "YYYY-MM-DDTHH:mm:ss",
        displayFormat: "MMM DD, YYYY, hh:mm:ss A",
        inputFormat: "MM/DD/YYYY hh:mm:ss A"
    },
    date: {
        valueFormat: "YYYY-MM-DD",
        displayFormat: "MMM DD, YYYY",
        inputFormat: "MM/DD/YYYY"
    },
    time: {
        valueFormatPrecision: "HH:mm:ss.SSS",
        valueFormat: "HH:mm:ss",
        displayFormat: "hh:mm:ss A",
        inputFormat: "hh:mm:ss A"
    }
}


export function formatNumberLocale() {
    const lng = i18next.resolvedLanguage ?? fallbackLng;
  
    const currentLang = allLangs.find((lang) => lang.value === lng);
  
    return { code: currentLang?.numberFormat.code, currency: currentLang?.numberFormat.currency };
  }
  