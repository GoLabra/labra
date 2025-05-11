
import { FieldScalarTypes } from "@/types/field-type-descriptor";
import { useMemo } from "react";
import { Field } from "@/lib/apollo/graphql.entities";
import { AdvancedFilter, AdvanedFilterOperator, AdvancedFilterProperty } from "@/core-features/dynamic-filter/filter";
import { useDialog } from "@/hooks/use-dialog";
import { SearchField } from "./content-manager-search";
import { containsFoldOperator, containsOperator, hasPrefixOperator, hasSufixOperator, eqStringFoldOperator, eqStringOperator, neqOperator, eqNumberOperator, neqNumberOperator, greaterThanNumberOperator, greaterThanOrEqualNumberOperator, lessThanNumberOperator, lessThanOrEqualNumberOperator, greaterThanDateTimeOperator, greaterThanOrEqualDateTimeOperator, lessThanDateTimeOperator, lessThanOrEqualDateTimeOperator, eqBooleanOperator } from "@/core-features/dynamic-filter/filter-operators";

export type FilterProperty2 = Partial<{
    [key in FieldScalarTypes]: {
        operators: AdvanedFilterOperator[];
    };
}>

const FieldTypeOperatorsMap: FilterProperty2 = {
    ShortText: {
        operators: [containsFoldOperator, containsOperator, hasPrefixOperator, hasSufixOperator, eqStringFoldOperator, eqStringOperator, neqOperator],
    },
    LongText: {
        operators: [containsFoldOperator, containsOperator, hasPrefixOperator, hasSufixOperator, eqStringFoldOperator, eqStringOperator, neqOperator],
    },
    Email: {
        operators: [containsFoldOperator, containsOperator, hasPrefixOperator, hasSufixOperator, eqStringFoldOperator, eqStringOperator, neqOperator],
    },
    RichText: {
        operators: [containsFoldOperator, containsOperator],
    },
    Integer: {
        operators: [eqNumberOperator, neqNumberOperator, greaterThanNumberOperator, greaterThanOrEqualNumberOperator, lessThanNumberOperator, lessThanOrEqualNumberOperator],
    },
    Decimal: {
        operators: [eqNumberOperator, neqNumberOperator, greaterThanNumberOperator, greaterThanOrEqualNumberOperator, lessThanNumberOperator, lessThanOrEqualNumberOperator],
    },
    Float: {
        operators: [eqNumberOperator, neqNumberOperator, greaterThanNumberOperator, greaterThanOrEqualNumberOperator, lessThanNumberOperator, lessThanOrEqualNumberOperator],
    },
    DateTime: {
        operators: [greaterThanDateTimeOperator, greaterThanOrEqualDateTimeOperator, lessThanDateTimeOperator, lessThanOrEqualDateTimeOperator],
    },
    Boolean: {
        operators: [ eqBooleanOperator ],
    },
    SingleChoice: {
        operators: [neqOperator, eqStringOperator],
    },
    MultipleChoice: {
        operators: [neqOperator, eqStringOperator],
    },
};


interface useContentManagerFilterDialogProps {
    fields: SearchField[];
}
export const useContentManagerFilterDialog = (props: useContentManagerFilterDialogProps) => {

    const { fields } = props;

    const dialog = useDialog();

    const filterProperties: AdvancedFilterProperty[] = useMemo((): AdvancedFilterProperty[] => {

        var result = fields.map(i => {

            var operators = FieldTypeOperatorsMap[i.type]?.operators ?? [];
            if (!operators.length) {
                return null;
            }

            return {
                label: i.label,
                name: i.name,
                operators: operators,
                selectOptions: i.selectOptions ?? []
            }
        }).filter(i => !!i);

        return result;
    }, [fields]);

    return {
        dialog,
        filterProperties
    }
}