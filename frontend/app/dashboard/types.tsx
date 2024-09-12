export type MutualFundInvestments = {
	investments: MutualFundInvestment[];
};

export type MutualFundInvestment = {
	scheme_id: number;
	scheme_name: string;
	current_nav: number;
	previous_day_nav: number;
	net_profit_loss_percentage: number;
	day_profit_loss_percentage: number;
	net_profit_loss: number;
	day_profit_loss: number;
	investments: MutualFundInvestmentDetails[];
};

export type MutualFundInvestmentDetails = {
	units: number;
	invested_at: string;
	invested_nav: number;
	current_value: number;
	invested_value: number;
	previous_day_value: number;
	net_profit_loss_percentage: number;
	day_profit_loss_percentage: number;
	net_profit_loss: number;
	day_profit_loss: number;
};
