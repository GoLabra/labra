import { NameCaptionEntity } from "@/types/entity";
import { makeVar, useReactiveVar } from "@apollo/client";
import { useCallback, useMemo } from "react";
import { appStatusVar, entitiesInCodeGenerationVar } from ".";
import { AppStatus } from "@/lib/apollo/graphql.entities";


export const useAppStatus = () => {
    const applicationStatus = useReactiveVar(appStatusVar);
    const entitiesInCodeGeneration = useReactiveVar(entitiesInCodeGenerationVar);

    const isEntityInCodeGeneration = useCallback((entityName: string) => {
        if (applicationStatus != AppStatus.Generating) {
            return false;
        }

        return entitiesInCodeGeneration.map(i => i.name).includes(entityName);
    }, [entitiesInCodeGeneration, applicationStatus]);

    const isCodeGenerating = useMemo(() => {
        return applicationStatus === AppStatus.Generating;
    }, [applicationStatus]);

    const isEntityFullyAvailable = useCallback((entityName: string) => {
        const isInCodeGenerating = isEntityInCodeGeneration(entityName);
        if(isInCodeGenerating){
            return false;
        }

        if(applicationStatus !== AppStatus.Up){
            return false;
        }

        return true;
    }, [isEntityInCodeGeneration, applicationStatus]);

    return useMemo(() => ({
        applicationStatus,
        entitiesInCodeGeneration,
        isEntityInCodeGeneration,
        isEntityFullyAvailable,
        isCodeGenerating
    }), [applicationStatus, entitiesInCodeGeneration, isEntityInCodeGeneration, isEntityFullyAvailable, isCodeGenerating]);
}