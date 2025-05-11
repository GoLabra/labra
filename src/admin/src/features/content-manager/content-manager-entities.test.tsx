import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { MockedProvider } from "@apollo/client/testing";
import ContentManagerEntities from "./content-manager-entities";
import { ThemeProvider } from "@mui/material";
import { theme } from "@/styles/theme";
import { GetEntitiesNameCaption } from "@/hooks/use-entities";



jest.mock("next/navigation", () => ({
    useParams: () => ({
        'entity-id': 'user',
    })
}));

jest.mock('@paralleldrive/cuid2', () => ({
    createId: jest.fn(() => 'mocked-id')
}));

jest.mock("@/config/CONST", () => ({
    COLOR: 'blue'
}));


function MockTheme({ children }: any) {

    return <ThemeProvider theme={theme}>{children}</ThemeProvider>;
}

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
                <ContentManagerEntities />
            </MockedProvider>
        </MockTheme>
    );

    const link = await screen.findByRole('link', { name: /artist/i });

    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', '/content-manager/artist');
});