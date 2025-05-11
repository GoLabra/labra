// import { fullEntitiesMapVar, GetEntitiesNameCaption, nameCaptionEntitiesVar } from "@/hooks/use-entities";
// import { NameCaptionEntity } from "@/types/entity";
// import { AppStatusEvent, appStatusVar, entitiesInCodeGenerationVar } from "../app-state";
// import { ApolloClient } from "@apollo/client";
// import { addNotification } from "@/lib/notifications/store";
// import { AppStatusApplicationUpMessage } from "@/types/centrifugo";

// export const handleApplicationShutdown = (
//     message: AppStatusApplicationUpMessage
// ) => {
    
//     appStatusVar(AppStatusEvent.SHUTDOWN);

//     // Add notification
//     addNotification({
//         message: "Application is shuting down/restarting",
//         type: 'info',
//     }, true);
// };