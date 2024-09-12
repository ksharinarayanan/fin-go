// declare type for mutual fund api response
// this is the response
export type MutualFundInvestment = {
	scheme_id: number;
	scheme_name: string;
	units: number;
	invested_at: string;
	current_nav: number;
	invested_nav: number;
	previous_day_nav: number;
	net_profit_loss_percentage: number;
	current_value: number;
	invested_value: number;
	previous_day_value: number;
	day_profit_loss_percentage: number;
	net_profit_loss: number;
	day_profit_loss: number;
};

export type MutualFundInvestmentsResponse = {
	mutual_funds: MutualFundInvestment[];
};
