import { Edge, Entity, Field, RelationType } from "@/lib/apollo/graphql.entities";
import { ReactNode, useCallback, useMemo } from "react";
import { FullEntity } from "@/types/entity";
import { z } from "zod";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { TextLongFormField } from "@/core-features/dynamic-form/form-fields/TextLongField";
import { RichTextFormField } from "@/core-features/dynamic-form/form-fields/RichTextField";
import { NumberFormField } from "@/core-features/dynamic-form/form-fields/NumberField";
import { dateTimeToString, dateTimeToStringPrecision, dateToString, stringToDate, stringToDateTime, stringToTime, timeToString, timeToStringPrecision } from "@/core-features/dynamic-form/value-convertor";
import { DateTimeFormField } from "@/core-features/dynamic-form/form-fields/DateTimeField";
import { TimeFormField } from "@/core-features/dynamic-form/form-fields/TimeField";
import { DateFormField } from "@/core-features/dynamic-form/form-fields/DateField";
import { BooleanSelectFormField } from "@/core-features/dynamic-form/form-fields/BooleanSelectField";
import { JSONFormField } from "@/core-features/dynamic-form/form-fields/JSONField";
import { SingleChoiceFormField } from "@/core-features/dynamic-form/form-fields/SingleChoice";
import dayjs, { Dayjs } from "dayjs";
import { LookupManyFIELDFormField } from "@/core-features/dynamic-form/form-fields/LookupManyFIELD";
import { LookupOneFIELDFormField, OptionDiffWrapper } from "@/core-features/dynamic-form/form-fields/LookupOneFIELD";
import { Options, Option } from "@/core-features/dynamic-form/form-field";
import { MultipleChoiceFormField } from "@/core-features/dynamic-form/form-fields/MultipleChoice";
import { isJsonOrNull, toJsonOrNull } from "@/lib/utils/is-json";
import { Upload } from "@mui/icons-material";
import { UploadFilesBaseField } from "@/core-features/dynamic-form/form-fields/UploadFilesBaseField";
import { FileDiffWrapper, FileFormField } from "@/core-features/dynamic-form/form-fields/FileField";
import { fileToBase64 } from "@/lib/utils/file-to-base64";
import { createId } from "@paralleldrive/cuid2";


type FieldDetails = {
	schema: z.ZodTypeAny;
	input: ReactNode;
	convertFromRawValue?: (value: any) => any;
}

const getShortText = (field: Field): FieldDetails => {

	const schema = (() => {
		if (field.required) {
			return z.preprocess(i => i ?? '',
				z.string().min(1, `${field.caption} is required`)
			);
		}
		return z.string().optional().nullish();
	})();

	return {
		schema,
		input: <TextShortFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} />
	}
}

const getLongText = (field: Field): FieldDetails => {

	const schema = (() => {
		if (field.required) {
			return z.preprocess(i => i ?? '',
				z.string().min(1, `${field.caption} is required`)
			);
		}
		return z.string().optional().nullish();
	})();

	return {
		schema,
		input: <TextLongFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} />
	}
}


const getRichText = (field: Field): FieldDetails => {

	const schema = (() => {
		if (field.required) {
			return z.preprocess(i => i ?? '',
				z.string().min(1, `${field.caption} is required`)
			);
		}
		return z.string().optional().nullish();
	})();

	return {
		schema,
		input: <RichTextFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} />
	}
}

const getEmail = (field: Field): FieldDetails => {
	const schema = (() => {
		if (field.required) {
			return z.preprocess(
				(val) => (typeof val === "string" ? val.trim() : val),
				z.string().min(1, `${field.caption} is required`).email('Invalid Email')
			);
		}
		return z.preprocess(
			(val) => (typeof val === "string" && val.trim() === "" ? undefined : val),
			z.string().email("Invalid Email").optional().nullish()
		);
	})();

	return {
		schema,
		input: <TextShortFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} />
	}
}

const getNumber = (field: Field): FieldDetails => {

	// 1) Normalize empty → null, otherwise pass through to z.coerce.number()
	let schema: z.ZodTypeAny = z.preprocess((value) => {
		if (value === null || value === undefined || value === "") {
			return null;
		}
		return value;
	}, z.coerce.number().nullable());

	// 2) If required, disallow null
	if (field.required) {
		schema = schema.refine((val) => val !== null, {
			message: `${field.caption} is required`,
		});
	}

	// Helper: only run test when there _is_ a number
	const isNullOr = (fn: (n: number) => boolean) =>
		(val: number | null) => val === null || fn(val);

	// 3) Min check (if configured)
	if (field.min != null) {
		const minVal = +field.min;
		schema = schema.refine(
			isNullOr((n) => n >= minVal),
			{ message: `${field.caption} must be ≥ ${minVal}` }
		);
	}

	// 4) Max check (if configured)
	if (field.max != null) {
		const maxVal = +field.max;
		schema = schema.refine(
			isNullOr((n) => n <= maxVal),
			{ message: `${field.caption} must be ≤ ${maxVal}` }
		);
	}


	return {
		schema,
		input: <NumberFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} min={field.min ?? undefined} max={field.max ?? undefined} />
	}
}


const getDateTime = (field: Field): FieldDetails => {

	// 1. Preprocess empty → null; strings → Dayjs
	let schema: z.ZodTypeAny = z.preprocess((value) => {
		if (value === null || value === undefined || value === "") {
			return null;
		}
		if (typeof value === "string") {
			// try to parse
			const dt = stringToDateTime(value);
			return dt ?? value;
		}
		return value;
	}, z.any().nullable());

	// 2. “Required?” check
	if (field.required) {
		schema = schema.refine((val) => val !== null, {
			message: `${field.caption} is required`,
		});
	}

	// Helper: only run a test if there’s actually a date
	const isNullOr = (fn: (d: Dayjs) => boolean) =>
		(val: Dayjs | null) => val === null || fn(val);

	let min: Dayjs | undefined = undefined;
	let max: Dayjs | undefined = undefined;

	// 3. Min-date check
	if (field.min) {
		min = stringToDateTime(field.min)!;
		schema = schema.refine(
			isNullOr((d) => !d.isBefore(min)),
			{ message: `${field.caption} must be ≥ ${field.min}` }
		);
	}

	// 4. Max-date check
	if (field.max) {
		max = stringToDateTime(field.max)!;
		schema = schema.refine(
			isNullOr((d) => !d.isAfter(max)),
			{ message: `${field.caption} must be ≤ ${field.max}` }
		);
	}

	schema = schema.transform((val) => dateTimeToStringPrecision(val));

	return {
		schema,
		//defaultValue: field.defaultValue ? stringToDateTime(field.defaultValue) : undefined,
		convertFromRawValue: (value) => {
			return value ? stringToDateTime(value) : undefined;
		},
		input: <DateTimeFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} min={min} max={max} />
	}
}


const getDate = (field: Field): FieldDetails => {
	// 1. Preprocess empty → null; strings → Dayjs
	let schema: z.ZodTypeAny = z.preprocess((value) => {
		if (value === null || value === undefined || value === "") {
			return null;
		}
		if (typeof value === "string") {
			// try to parse
			const dt = stringToDate(value);
			return dt ?? value;
		}
		return value;
	}, z.any().nullable());

	// 2. “Required?” check
	if (field.required) {
		schema = schema.refine((val) => val !== null, {
			message: `${field.caption} is required`,
		});
	}

	// Helper: only run a test if there’s actually a date
	const isNullOr = (fn: (d: Dayjs) => boolean) =>
		(val: Dayjs | null) => val === null || fn(val);

	let min: Dayjs | undefined = undefined;
	let max: Dayjs | undefined = undefined;

	// 3. Min-date check
	if (field.min) {
		min = stringToDate(field.min)!;
		schema = schema.refine(
			isNullOr((d) => !d.isBefore(min)),
			{ message: `${field.caption} must be ≥ ${field.min}` }
		);
	}

	// 4. Max-date check
	if (field.max) {
		max = stringToDate(field.max)!;
		schema = schema.refine(
			isNullOr((d) => !d.isAfter(max)),
			{ message: `${field.caption} must be ≤ ${field.max}` }
		);
	}

	schema = schema.transform((val) => dateToString(val));

	return {
		schema,
		//defaultValue: field.defaultValue ? stringToDate(field.defaultValue) : undefined,
		convertFromRawValue: (value) => {
			return value ? stringToDate(value) : undefined;
		},
		input: <DateFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} min={min} max={max} />
	}
}


const getTime = (field: Field): FieldDetails => {
	// 1. Preprocess empty → null; strings → Dayjs
	let schema: z.ZodTypeAny = z.preprocess((value) => {
		if (value === null || value === undefined || value === "") {
			return null;
		}
		if (typeof value === "string") {
			// try to parse
			const dt = stringToTime(value);
			return dt ?? value;
		}
		return value;
	}, z.any().nullable());

	// 2. “Required?” check
	if (field.required) {
		schema = schema.refine((val) => val !== null, {
			message: `${field.caption} is required`,
		});
	}

	// Helper: only run a test if there’s actually a date
	const isNullOr = (fn: (d: Dayjs) => boolean) =>
		(val: Dayjs | null) => val === null || fn(val);

	let min: Dayjs | undefined = undefined;
	let max: Dayjs | undefined = undefined;

	// 3. Min-date check
	if (field.min) {
		min = stringToTime(field.min)!;
		schema = schema.refine(
			isNullOr((d) => !d.isBefore(min)),
			{ message: `${field.caption} must be ≥ ${field.min}` }
		);
	}

	// 4. Max-date check
	if (field.max) {
		max = stringToTime(field.max)!;
		schema = schema.refine(
			isNullOr((d) => !d.isAfter(max)),
			{ message: `${field.caption} must be ≤ ${field.max}` }
		);
	}

	schema = schema.transform((val) => timeToStringPrecision(val));

	return {
		schema,
		//defaultValue: field.defaultValue ? stringToTime(field.defaultValue) : undefined,
		convertFromRawValue: (value) => {
			return value ? stringToTime(value) : undefined;
		},
		input: <TimeFormField key={field.name} name={field.name} label={field.caption} required={field.required ?? false} min={min} max={max} />
	}
}

const getBoolean = (field: Field): FieldDetails => {
	const baseSchema = z.boolean();

	let schema;
	if (field.required) {
		schema = baseSchema.refine(
			(data) => data !== undefined,
			{ message: `${field.caption} is required` }
		);
	} else {
		schema = baseSchema.optional();
	}
	return {
		schema,
		//defaultValue: field.defaultValue ? !!field.defaultValue : undefined,
		convertFromRawValue: (value) => {
			return value ? !!value : undefined;
		},
		input: <BooleanSelectFormField key={field.name} name={field.name} label={field.caption} />
	}
}


const getJson = (field: Field): FieldDetails => {

	let schema: z.ZodTypeAny = (() => {
		if (field.required) {
			return z.preprocess(i => i ?? '',
				z.string().min(1, `${field.caption} is required`)
			);
		}
		return z.string().transform((i) => i?.trim() || undefined).optional().nullish();
	})();


	// 2) If required, ban empty/null
	if (field.required) {
		schema = schema.refine(
			(val) => val !== undefined && val !== null,
			{ message: `${field.caption} is required` }
		);
	}

	// 3) JSON‐validity only if there is a value
	schema = schema.refine(
		(val) => val === undefined || val === null || isJsonOrNull(val),
		{ message: `${field.caption} must be valid JSON` }
	).transform((val) => toJsonOrNull(val) ? JSON.parse(val) : undefined);

	return {
		schema,
		convertFromRawValue: (value) => {
			if (typeof value === 'string') {
				return value;
			}
			return value ? JSON.stringify(value) : undefined;
		},
		input: <JSONFormField key={field.name} name={field.name} label={field.caption} />
	}
}

const getSingleChoice = (field: Field): FieldDetails => {

	let schema = field.required
		? z.any().refine((val) => !!val, { message: `${field.caption} is required` })
		: z.string().optional().nullable();

	const options = field.acceptedValues?.filter(i => i)
		.map((i: string | null) => ({
			label: i!,
			value: i!
		})) ?? [];

	return {
		schema,
		input: <SingleChoiceFormField key={field.name} name={field.name} label={field.caption} options={options} />
	}
}


const getMultiChoice = (field: Field): FieldDetails => {

	// let schema = field.required
	//     ? z.any().refine((val) => !!val && val.length, { message: `${field.caption} is required` })
	//     : z.string().optional().nullable();


	const schema = (() => {
		if (field.required) {
			return z.preprocess(i => i ?? [],
				z.array(z.string()).min(1, `${field.caption} is required`)
			);
		}
		return z.array(z.string()).optional().nullish();
	})();

	const options = field.acceptedValues?.filter(i => i)
		.map((i: string | null) => ({
			label: i!,
			value: i!
		})) ?? [];

	return {
		schema,
		//defaultValue: field.defaultValue ? JSON.parse(field.defaultValue) : undefined,
		convertFromRawValue: (value) => {
			if (Array.isArray(value)) {
				return value;
			}
			return value ? JSON.parse(value) : undefined;
		},
		input: <MultipleChoiceFormField key={field.name} name={field.name} label={field.caption} options={options} />
	}
}

const getUploadOneFile = (entityName: string, edge: Edge): FieldDetails => {

	let schema: z.ZodTypeAny = edge.required
		? z.any().refine((val): val is Record<string, unknown> => !val,
			{ message: `${edge.caption} is required` }
		)
		: z.any().optional().nullable();

	schema = schema.transform(async (val?: FileDiffWrapper[]) => {

		if (!val) {
			return undefined;
		}

		if (!val.length) {
			return undefined;
		}

		const file = val[0];

		if(!file){
			return undefined;
		}

		switch (file.status) {
			case 'create':
				return {
					create: {
						name: (file.file as File).name,
						content: await fileToBase64(file.file as File),

						//!!!! remove these
						storageFileName: createId(),
						size: 1,
					}
				}
			case 'delete':
				return {
					delete: true
				}
		}

	});

	return {
		schema,
		convertFromRawValue: (val) => undefined,
		input: <FileFormField key={edge.name} name={edge.name} label={edge.caption} entityName={entityName} edge={edge} maxFiles={1} />
	}
}

const getUploadManyFile = (entityName: string, edge: Edge): FieldDetails => { 

	let schema: z.ZodTypeAny = edge.required
		? z.any().refine((val): val is Record<string, unknown> => !val,
			{ message: `${edge.caption} is required` }
		)
		: z.any().optional().nullable();

	schema = schema.transform(async (val?: FileDiffWrapper[]) => {

		if (!val) {
			return undefined;
		}

		const create = val.filter(i => i.status === 'create' && i.file instanceof File);
		const remove = val.filter(i => i.status === 'delete');

		const result = {
			...(create.length && {
				create: await Promise.all(val.filter(i => i.status === 'create')
					.filter((i) => i.file instanceof File)
					.map(async (i: FileDiffWrapper) => ({
						name: (i.file as File).name,
						content: await fileToBase64(i.file as File),

						//!!!! remove these
						storageFileName: createId(),
						size: 1,

					})))
			}),
			...(remove.length && {
				delete: val.filter(i => i.status === 'delete')
					.map(i => ({
						id: i.id
					}))
			})
		}
		console.log(result);
		return result;
	});

	return {
		schema,
		convertFromRawValue: (val) => undefined,
		input: <FileFormField key={edge.name} name={edge.name} label={edge.caption} entityName={entityName} edge={edge} />
	}
}


const getRelationOne = (entityName: string, edge: Edge): FieldDetails => {

	let schema: z.ZodTypeAny = edge.required
		? z.any().refine((val): val is Record<string, unknown> => !val,
			{ message: `${edge.caption} is required` }
		)
		: z.any().optional().nullable();

	schema = schema.transform((val?: OptionDiffWrapper) => {

		if (!val) {
			return undefined;
		}

		if (!val.status) {
			return undefined;
		}

		switch (val.status) {
			case 'create':
				return {
					create: val.value
				}
			case 'connect':
				return {
					connect: {
						id: val.value
					}
				}
			case 'unset':
				return {
					unset: true
				}
			case 'delete':
				return {
					delete: true
				}
		}

	});

	return {
		schema,
		convertFromRawValue: (val) => undefined,
		input: <LookupOneFIELDFormField key={edge.name} name={edge.name} label={edge.caption} entityName={entityName} edge={edge} />
	}
}

const getRelationMany = (entityName: string, edge: Edge): FieldDetails => {
	let schema: z.ZodTypeAny = edge.required
		? z.any().refine((val): val is Record<string, unknown> => !val,
			{ message: `${edge.caption} is required` }
		)
		: z.any().optional().nullable();

	schema = schema.transform((val?: OptionDiffWrapper[]) => {

		if (!val || !val.length) {
			return undefined;
		}

		const connect = val.filter(i => i.status === 'connect'); 
		const create = val.filter(i => i.status === 'create');
		const disconnect = val.filter(i => i.status === 'disconnect');
		const remove = val.filter(i => i.status === 'delete');

		return {
			...(connect.length && {
				connect: connect.map(i => ({
					id: i.id
				})),
			}),
			...(create.length && {
				create: create.map(i => i.value),
			}),
			...(disconnect.length && {
				disconnect: disconnect.map(i => ({
					id: i.id
				}))
			}),
			...(remove.length && {
				delete: remove.map(i => ({
					id: i.id
				}))
			})
		}
	});

	return {
		schema,
		convertFromRawValue: (val) => undefined,
		input: <LookupManyFIELDFormField key={edge.name} name={edge.name} label={edge.caption} entityName={entityName} edge={edge} />
	}
}

const getFormField = (field: Field): FieldDetails => {
	switch (field.type) {
		case "ShortText": return getShortText(field);
		case "LongText": return getLongText(field);
		case "RichText": return getRichText(field);
		case "Email": return getEmail(field);
		case "Integer": return getNumber(field);
		case "Decimal": return getNumber(field);
		case "Float": return getNumber(field);
		case "DateTime": return getDateTime(field);
		case "Date": return getDate(field);
		case "Time": return getTime(field);
		case "Media": return getShortText(field);
		case "Boolean": return getBoolean(field);
		case "Json": return getJson(field);
		case "SingleChoice": return getSingleChoice(field);
		case "MultipleChoice": return getMultiChoice(field);
		default: return getShortText(field);
	}
}

const getFormEdge = (entityName: string, edge: Edge): FieldDetails => {

	const isFile = edge.relatedEntity.caption == 'File';

	switch (edge.relationType) {
		case RelationType.One:
		case RelationType.OneToOne:
		case RelationType.OneToMany:

			if (isFile) {
				return getUploadOneFile(entityName, edge);
			}

			return getRelationOne(entityName, edge);

		case RelationType.Many:
		case RelationType.ManyToOne:
		case RelationType.ManyToMany:

			if (isFile) {
				return getUploadManyFile(entityName, edge);
			}

			return getRelationMany(entityName, edge);

		default:
			return {
				schema: z.any().nullable(),
				input: <>Unknown Relation</>
			}
	}

}


export const useContentManagerFormSchema = (entity: FullEntity | null) => {

	const fieldsDescriptors = useMemo(() => {
		if (!entity) {
			return [];
		}

		const fields = entity.fields!
			.filter(i => i.name != 'id') // skip id
			.filter(i => i.name != 'createdAt' && i.name != 'updatedAt') // TODO: skip createdAt, updateAt from form schema
			.map((field) => ({
				...getFormField(field),
				name: field.name,
				field
			}));

		const edges = entity.edges!
			.filter(i => i.name.startsWith('ref') == false) //TODO: ref fields will be removed from BE
			.filter(i => i.name != 'createdBy' && i.name != 'updatedBy') // skip createdBy, updatedBy from form schema
			.map((edge) => ({
				...getFormEdge(entity.name, edge),
				name: edge.name,
				edge
			}));

		return [
			...fields,
			...edges
		]

	}, [entity]);

	const schema = useMemo(() => {
		const fieldsSchema = fieldsDescriptors.reduce((acc: any, field) => {
			acc = {
				...acc,
				[field.name]: field.schema
			}
			return acc;
		}, {} as Record<string, z.ZodTypeAny>);

		return z.object(fieldsSchema);
	}, [fieldsDescriptors]);

	const fields = useMemo(() => {
		return fieldsDescriptors.map(field => field.input);
	}, [fieldsDescriptors]);

	const schemaDefaultValue = useMemo(() => {
		return entity?.fields.reduce((acc: any, field) => {
			if (field.defaultValue) {
				acc[field.name] = field.defaultValue;
			}
			return acc;
		}, {});
	}, [entity?.fields]);

	const convertFromRawValue = useCallback((defaultValue: any) => {
		if (!defaultValue) {
			return defaultValue;
		}

		const keys = Object.keys(defaultValue);
		const result = keys.reduce((acc: any, key) => {
			const feldDescription = fieldsDescriptors.find(i => i.name == key);
			const value = feldDescription?.convertFromRawValue ? feldDescription?.convertFromRawValue?.(defaultValue[key]) : defaultValue[key];

			if (value == undefined || value == null) {
				return acc;
			}

			acc[key] = value;
			return acc;
		}, {});

		return result;
	}, []);

	return useMemo(() => ({
		schema,
		fields,
		schemaDefaultValue,
		convertFromRawValue
	}), [schema, fields]);
}

useContentManagerFormSchema.displayName = 'useContentManagerFormSchema';


export const useGenericFiedFormSchema = () => {

}

export const useGenericEdgeFormSchema = (entity: FullEntity | null, edgeName: string) => {

	const edge = useMemo(() => {
		if (!entity) {
			return null;
		}
		return entity.edges?.find(i => i.name === edgeName);
	}, [entity, edgeName]);

	const edgeDescriptor = useMemo(() => {
		if (!entity) {
			return null;
		}

		if (!edge) {
			return null;
		}

		const result = {
			...getFormEdge(entity.name, edge),
			name: edge.name,
			edge
		};
		return result;
	}, [entity, edge]);

	const convertFromRawValue = useCallback((defaultValue: any) => {

		if (!edge) {
			return null;
		}

		const value = edgeDescriptor?.convertFromRawValue ? edgeDescriptor?.convertFromRawValue?.(defaultValue) : defaultValue;
		return value;
	}, [edgeDescriptor, edge]);

	const field = useMemo((): ReactNode => {
		if (!edgeDescriptor) {
			return <></>;
		}

		return edgeDescriptor.input;
	}, [edgeDescriptor]);

	const schema = useMemo(() => {
		if (!edgeDescriptor) {
			return z.any();
		}

		return edgeDescriptor.schema;
	}, [edgeDescriptor]);

	return useMemo(() => ({
		edgeDescriptor,
		convertFromRawValue,
		field,
		schema
	}), [edgeDescriptor, convertFromRawValue, field, schema]);
}