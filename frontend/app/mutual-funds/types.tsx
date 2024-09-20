export type Scheme = {
	ID: number;
	SchemeName: string;
};

export type SchemeData = {
	schemeCode: number;
	schemeName: string;
	currentNav: number;
	date: string;
};

export type AddInvestmentResponse = {
	message: string;
};
