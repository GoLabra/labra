import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { Box } from "@mui/material";
import { ContentManagerEntryGeneric } from "./content-manager-entry-generic";
import { MockedProvider } from "@apollo/client/testing";
import * as dialogContextModule from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import * as useEntitiesModule from "@/hooks/use-entities";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { ForwardRefExoticComponent, RefAttributes } from "react";
import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { FullEntity } from "@/types/entity";
import { EntityOwner } from "@/lib/apollo/graphql.entities";

jest.mock("@/core-features/dynamic-dialog/src/use-my-dialog-context", () => ({
	useMyDialogContext: jest.fn()
}));

jest.mock("@/hooks/use-entities", () => ({
	useFullEntity: jest.fn()
}));

jest.mock('@paralleldrive/cuid2', () => ({
	createId: jest.fn(() => 'mocked-id')
}));

describe("Generic Entry Form - Create New/Edit Mode", () => {

	it("Do not show ID label and the timeline on create new entry", async () => {
		const mockContextValue = {
			setLocalState: jest.fn(),
			upperResults: {},
			localState: {},
			openMode: FormOpenMode.New,
			editId: null,
			close: jest.fn(),
			closeWithResults: jest.fn(),
			addPopup: jest.fn(),
			addPopupAsFirst: jest.fn()
		};

		const mockFullEntity = {
			name: 'TestEntity',
			caption: 'TestEntity',
			owner: EntityOwner.User,
			displayField: {
				__typename: 'Field',
				name: 'name',
				caption: 'Name',
				type: 'ShortText',
				required: false,
				unique: false,
			},
			fields: [],
			edges: [],
			loading: false,
		};

		(dialogContextModule.useMyDialogContext as jest.Mock).mockReturnValue(mockContextValue);
		(useEntitiesModule.useFullEntity as jest.Mock).mockReturnValue(mockFullEntity);

		const { container } = render(
			<MockedProvider>
				<ContentManagerEntryGeneric entityName="TestEntity" defaultValue={{ }} />
			</MockedProvider>
		);

		const elementWithId = container.querySelector('#test-entry-id-777');
		expect(elementWithId).not.toBeInTheDocument();

		const timelineElement = container.querySelector('.MuiTimeline-root');
		expect(timelineElement).not.toBeInTheDocument();

	});

	it("Show ID label and the timeline on edit existing entry", async () => {
		const mockContextValue = {
			setLocalState: jest.fn(),
			upperResults: {},
			localState: {},
			openMode: FormOpenMode.Edit,
			editId: 'test-entry-id-777',
			close: jest.fn(),
			closeWithResults: jest.fn(),
			addPopup: jest.fn(),
			addPopupAsFirst: jest.fn()
		};

		const mockFullEntity = {
			name: 'user',
			caption: 'User',
			owner: EntityOwner.User,
			displayField: {
				__typename: 'Field',
				name: 'name',
				caption: 'Name',
				type: 'ShortText',
				required: false,
				unique: false,
			},
			fields: [],
			edges: [],
			loading: false,
		};

		(dialogContextModule.useMyDialogContext as jest.Mock).mockReturnValue(mockContextValue);
		(useEntitiesModule.useFullEntity as jest.Mock).mockReturnValue(mockFullEntity);

		const { container } = render(
			<MockedProvider>
				<ContentManagerEntryGeneric entityName="TestEntity" defaultValue={{ }} />
			</MockedProvider>
		);

		const elementWithId = container.querySelector('#test-entry-id-777');
		expect(elementWithId).toBeInTheDocument();

		const timelineElement = container.querySelector('.MuiTimeline-root');
		expect(timelineElement).toBeInTheDocument();

	});

});