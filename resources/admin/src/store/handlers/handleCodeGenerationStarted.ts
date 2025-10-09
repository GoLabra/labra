// import { fullEntitiesMapVar, GetEntitiesNameCaption, nameCaptionEntitiesVar } from "@/hooks/use-entities";
// import { NameCaptionEntity } from "@/types/entity";
// import { AppStatusEvent, appStatusVar, entitiesInCodeGenerationVar } from "../app-state";
// import { ApolloClient } from "@apollo/client";
// import { addNotification } from "@/lib/notifications/store";
// import { AppStatusEventCodeGenerationStartedMessage } from "@/types/centrifugo";

// export const handleCodeGenerationStarted = (
//     message: AppStatusEventCodeGenerationStartedMessage,
//     client: ApolloClient<any>
// ) => {
//     const nameCaptionEntities: NameCaptionEntity[] = message.entities?.map(entity => ({
//         name: entity.Name,
//         caption: entity.Caption,
//     })) ?? [];

//     // Clean up full entities cache
//     const unchangedEntities = Object.fromEntries(
//         Object.entries(fullEntitiesMapVar()).filter(
//             ([key]) => !nameCaptionEntities.some(n => n.name === key)
//         )
//     );
//     fullEntitiesMapVar(unchangedEntities);

//     // Update code generation state
//     appStatusVar(AppStatusEvent.CODE_GENERATION_STARTED);
//     entitiesInCodeGenerationVar(nameCaptionEntities);

//     // Update Apollo cache with new entities
//     client.writeQuery({
//         query: GetEntitiesNameCaption,
//         data: {
//             entities: [
//                 ...nameCaptionEntitiesVar(),
//                 ...nameCaptionEntities.filter(entity => !nameCaptionEntitiesVar().some(n => n.name === entity.name)),
//             ].sort((a, b) => a.name.localeCompare(b.name)),
//         },
//     });

//     // Add notification
//     addNotification({
//         message: "Please continue using the application. You'll be able to modify the entity once we're back into full color",
//         type: 'info'
//     }, true);
// };