// import { fullEntitiesMapVar, GetEntitiesNameCaption, nameCaptionEntitiesVar } from "@/hooks/use-entities";
// import { NameCaptionEntity } from "@/types/entity";
// import { AppStatusEvent, appStatusVar, entitiesInCodeGenerationVar } from "../app-state";
// import { ApolloClient } from "@apollo/client";
// import { addNotification } from "@/lib/notifications/store";
// import { AppStatusApplicationUpMessage, AppStatusEventCodeGenerationFailedMessage } from "@/types/centrifugo";

// export const handleCodeGenerationFailed = (
//     message: AppStatusEventCodeGenerationFailedMessage
// ) => {
    
//     appStatusVar(AppStatusEvent.CODE_GENERATION_FAILED);

//     // Add notification
//     addNotification({
//         message: "Code generation failed. All changes will be reverted",
//         type: 'info'
//     }, true);
// };