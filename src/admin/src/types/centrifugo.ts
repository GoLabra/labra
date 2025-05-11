// import { AppStatus } from "@/store/app-state";


// /**************
//  * APP_STATUS
//  *************/
// export interface AppStatusEventMessage {
//     event: AppStatus;
//     timestamp: string;
// }

// export interface AppStatusApplicationUpMessage extends AppStatusEventMessage {
//     event: typeof AppStatus['UP'] ;
// }

// export interface AppStatusEventCodeGenerationStartedMessage extends AppStatusEventMessage {
//     event: typeof AppStatus['CODE_GENERATION_STARTED'] ;
//     entities?: {
//         Name: string;
//         EntName: string;
//         Caption: string;
//         Owner: string;
//         DisplayFieldName: string;
//     }[];
// }

// export interface AppStatusEventCodeGenerationCompletedMessage extends AppStatusEventMessage {
//     event: typeof AppStatus['CODE_GENERATION_COMPLETED'] ;
// }

// export interface AppStatusEventCodeGenerationFailedMessage extends AppStatusEventMessage {
//     event: typeof AppStatus['CODE_GENERATION_FAILED'] ;
// }

// export interface AppStatusEventCodeGenerationRevertedMessage extends AppStatusEventMessage {
//     event: typeof AppStatus['CODE_GENERATION_REVERTED'] ;
// }

// export interface AppStatusEventShutdownMessage extends AppStatusEventMessage {
//     event: typeof AppStatus['SHUTDOWN'] ;
// }

// /**************
//  * USER-MESSAGES
//  *************/
// export interface UserMessageMessage {
//     message: string;
//     severity?: 'info' | 'success' | 'warning' | 'error';
// }