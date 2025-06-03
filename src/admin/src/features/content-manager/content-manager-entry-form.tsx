import { forwardRef } from "react"; 
import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { ContentManagerEntryGeneric } from "./form-entry/content-manager-entry-generic";
import { ContentManagerEntryRole } from "./form-entry/content-manager-entry-role";

export interface ContentManagerEntryDialogContentProps {
	entityName: string,
	defaultValue: any
}
export const ContentManagerEntryDialogContent = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

	switch (props.entityName) {
		case 'role':
			return <ContentManagerEntryRole key="role" {...props} ref={ref} />;
		default:
			return <ContentManagerEntryGeneric key="generic" {...props} ref={ref} />;
	}
});
ContentManagerEntryDialogContent.displayName = 'ContentManagerEntryDialogContent';

