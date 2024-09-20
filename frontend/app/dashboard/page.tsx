"use client";

import React, { useEffect, useState } from "react";
import { MutualFundInvestment, MutualFundInvestments } from "./types";
import MutualFundAccordion from "./components";
import { Accordion } from "@/components/ui/accordion";
import { round } from "../../utils/utils";

const Dashboard = () => {
	const [mfInvestments, setMfInvestments] = useState<MutualFundInvestment[]>(
		[]
	);
	const [isLoading, setIsLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const deriveValues = (mf: MutualFundInvestment[]) => {
		for (const investment of mf) {
			let currentTotalValue: number = 0,
				investedTotalValue: number = 0,
				previousDayTotalValue: number = 0;

			for (const investmentDetail of investment.investments) {
				currentTotalValue += investmentDetail.current_value;
				investedTotalValue += investmentDetail.invested_value;
				previousDayTotalValue += investmentDetail.previous_day_value;
			}

			investment.net_profit_loss = round(
				currentTotalValue - investedTotalValue
			);
			investment.day_profit_loss = round(
				currentTotalValue - previousDayTotalValue
			);

			investment.net_profit_loss_percentage = round(
				(investment.net_profit_loss / investedTotalValue) * 100
			);
			investment.day_profit_loss_percentage = round(
				(investment.day_profit_loss / previousDayTotalValue) * 100
			);
		}
	};

	useEffect(() => {
		const fetchData = async () => {
			try {
				const response = await fetch(
					"http://localhost:8080/api/mf/investments"
				);
				if (!response.ok) {
					throw new Error("Network response was not ok");
				}
				const data: MutualFundInvestments = await response.json();
				deriveValues(data.investments);
				setMfInvestments(data.investments);
			} catch (error) {
				setError("Failed to fetch mutual fund data");
				console.error("Error fetching mutual fund data:", error);
			} finally {
				setIsLoading(false);
			}
		};

		fetchData();
	}, []);

	if (isLoading) {
		return <div className="text-center py-10">Loading...</div>;
	}

	if (error) {
		return <div className="text-center py-10 text-red-500">{error}</div>;
	}

	return (
		<div className="container mx-auto p-4">
			<h1 className="text-3xl font-bold mb-6">Mutual Fund Investments</h1>
			<Accordion type="single" collapsible className="w-full">
				{mfInvestments.map((mf: MutualFundInvestment) => (
					<MutualFundAccordion key={mf.scheme_id} mf={mf} />
				))}
			</Accordion>
		</div>
	);
};

export default Dashboard;
