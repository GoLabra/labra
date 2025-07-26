import "@testing-library/jest-dom";
import { render, screen, within } from "@testing-library/react";
import userEvent from '@testing-library/user-event';
import { MockedProvider } from "@apollo/client/testing";
import { ThemeProvider } from "@mui/material";
import EntityTypeBuilderEntities from "./entity-type-designer-entities";
import { theme } from "@/styles/theme";
import { GetEntitiesNameCaption } from "@/hooks/use-entities";

function MockTheme({ children }: any) {

    return <ThemeProvider theme={theme}>{children}</ThemeProvider>;
}

jest.mock('@paralleldrive/cuid2', () => ({
    createId: jest.fn(() => 'mocked-id')
}));

jest.mock("@/config/CONST", () => ({
    COLOR: 'blue',
	ENTITY_SYSTEM_KEYWORDS: []
}));

// Mock useRouter:
jest.mock("next/navigation", () => ({
    useRouter() {
        return {
            prefetch: () => null,
            push: (path: string) => null,
        };
    },
    useParams: () => ({
        'entity-id': 'user',
    }),
    usePathname: jest.fn(() => "localhost:3000/your-mocked-pathname"),
}));

describe("Load graphql data", () => {

    const mocks = [
        {
            request: {
                query: GetEntitiesNameCaption,
            },
            result: {
                "data": {
                    "entities": [
                        {
                            "name": "artist",
                            "caption": "Artist",
                            "displayField": {
                                "name": "name",
                                "caption": "Name",
                                "type": "ShortText",
                                "required": false,
                                "unique": false,
                                "defaultValue": null,
                                "min": null,
                                "max": null,
                                "private": null,
                                "acceptedValues": null,
                                "__typename": "Field"
                            },
                            "__typename": "Entity"
                        }
                    ]
                }
            }
        }
    ];

    it("renders the entities list", async () => {
        render(
            <MockTheme>
                <MockedProvider mocks={mocks} addTypename={false}>
                    <EntityTypeBuilderEntities />
                </MockedProvider>
            </MockTheme>
        );

        const link = await screen.findByRole('link', { name: /artist/i });

        expect(link).toBeInTheDocument();
        expect(link).toHaveAttribute('href', '/entity-type-designer/artist');
    });


    it("renders the add entity button", async () => {
        render(
            <MockTheme>
                <MockedProvider mocks={mocks} addTypename={false}>
                    <EntityTypeBuilderEntities />
                </MockedProvider>
            </MockTheme>
        );

        await screen.findByRole('button', { name: /add new entity/i });

        const button = screen.getByRole('button', { name: /add new entity/i });
        expect(button).toBeInTheDocument();
    });
});



describe("Handles the full add new entity flow", () => {

    const mocks = [
        {
            request: {
                query: GetEntitiesNameCaption,
            },
            result: {
                "data": {
                    "entities": []
                }
            }
        }
    ];

    it("renders the add entity button", async () => {
        render(
            <MockTheme>
                <MockedProvider mocks={mocks} addTypename={false}>
                    <EntityTypeBuilderEntities />
                </MockedProvider>
            </MockTheme>
        );

        await screen.findByRole('button', { name: /add new entity/i });
        
        const button = screen.getByRole('button', { name: /add new entity/i });
        
        await userEvent.click(button);

        
        const dialog = screen.getByRole('dialog');
        
        expect(dialog).toBeInTheDocument();
        
        await userEvent.type(screen.getByRole('textbox', { name: /caption/i }), 'My new entity');
        await userEvent.click(screen.getByRole('button', { name: /save/i }));

        const link = await screen.findByRole('link', { name: /My new entity/i });

        expect(link).toBeInTheDocument();
        expect(link).toHaveAttribute('href', '/entity-type-designer/mocked-id');
    })
});