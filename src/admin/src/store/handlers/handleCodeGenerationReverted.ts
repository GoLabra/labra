// import { fullEntitiesMapVar, GetEntitiesNameCaption } from "@/hooks/use-entities";
// import { NameCaptionEntity } from "@/types/entity";
// import { AppStatusEvent, appStatusVar, entitiesInCodeGenerationVar } from "../app-state";
// import { ApolloClient } from "@apollo/client";
// import { addNotification } from "@/lib/notifications/store";
// import { AppStatusEventCodeGenerationRevertedMessage } from "@/types/centrifugo";

// export const handleCodeGenerationReverted = (
//     message: AppStatusEventCodeGenerationRevertedMessage
// ) => {

//     appStatusVar(AppStatusEvent.CODE_GENERATION_REVERTED);

//     addNotification({
//         message: "Code generation is reverted",
//         type: 'success'
//     });
// };