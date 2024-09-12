import { MutualFundInvestment } from "./types";

export default function MutualFundCard({ mf }: { mf: MutualFundInvestment }) {
	return (
		<div className="flex flex-col gap-2 p-4 border border-gray-300 rounded-md">
			<h2 className="text-lg font-bold">{mf.scheme_name}</h2>
			<div className="grid grid-cols-2 gap-2 text-sm">
				<div>
					<p>Units: {mf.units.toFixed(2)}</p>
					<p>
						Invested on:{" "}
						{new Date(mf.invested_at).toLocaleDateString()}
					</p>
				</div>
				<div>
					<p>Current NAV: &#8377;{mf.current_nav.toFixed(2)}</p>
					<p>Invested NAV: &#8377;{mf.invested_nav.toFixed(2)}</p>
				</div>
			</div>
			<div className="grid grid-cols-2 gap-2 text-sm">
				<div>
					<p>Current Value: &#8377;{mf.current_value.toFixed(2)}</p>
					<p>Invested Value: &#8377;{mf.invested_value.toFixed(2)}</p>
				</div>
				<div>
					<p>
						Day Profit/Loss: &#8377;{mf.day_profit_loss.toFixed(2)}{" "}
						({mf.day_profit_loss_percentage.toFixed(2)}%)
					</p>
					<p>
						Net Profit/Loss: &#8377;{mf.net_profit_loss.toFixed(2)}{" "}
						({mf.net_profit_loss_percentage.toFixed(2)}%)
					</p>
				</div>
			</div>
		</div>
	);
}
