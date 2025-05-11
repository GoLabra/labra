import { AppStatus } from "@/lib/apollo/graphql.entities";
import { NameCaptionEntity } from "@/types/entity";
import { makeVar } from "@apollo/client";


export const appStatusVar = makeVar<AppStatus>(AppStatus.Up);
export const entitiesInCodeGenerationVar = makeVar<NameCaptionEntity[]>([]);
