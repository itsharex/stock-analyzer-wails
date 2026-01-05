export interface GlossaryItem {
  term: string;
  explanation: string;
  simpleExplanation: string;
}

export const STOCK_GLOSSARY: Record<string, GlossaryItem> = {
  "上升趋势线": {
    term: "上升趋势线",
    explanation: "连接股价波动过程中的低点而画出的直线，代表股价整体重心向上移。",
    simpleExplanation: "就像爬坡时的脚踏板，只要股价踩在上面，就说明还在往上涨。"
  },
  "短期强支撑位": {
    term: "短期强支撑位",
    explanation: "在短期内，股价下跌到某一价格水平时，买盘力量增强，阻止股价进一步下跌。",
    simpleExplanation: "股价的‘地板’，跌到这里通常会有很多人想买，不容易跌破。"
  },
  "前高压制位": {
    term: "前高压制位",
    explanation: "股价上涨到前期高点附近时，由于前期套牢盘或获利盘的抛售压力，导致股价难以突破。",
    simpleExplanation: "股价的‘天花板’，涨到之前的高点时，很多人想卖出解套，所以很难涨过去。"
  },
  "金叉": {
    term: "金叉",
    explanation: "短期指标线由下向上穿过长期指标线，通常被视为买入信号。",
    simpleExplanation: "快线超过了慢线，说明最近涨势变猛了，是个好兆头。"
  },
  "死叉": {
    term: "死叉",
    explanation: "短期指标线由上向下穿过长期指标线，通常被视为卖出信号。",
    simpleExplanation: "快线跌破了慢线，说明最近跌势变快了，要小心。"
  },
  "MACD": {
    term: "MACD",
    explanation: "指数平滑异同移动平均线，利用短期和长期均线的聚合与分离，研判买卖时机。",
    simpleExplanation: "股票的‘体温计’，用来判断股票现在是热（涨）还是冷（跌）。"
  },
  "KDJ": {
    term: "KDJ",
    explanation: "随机指标，通过计算最高价、最低价和收盘价，反映价格走势的强弱和超买超卖。",
    simpleExplanation: "股票的‘情绪指标’，看看大家现在是买疯了（超买）还是卖惨了（超卖）。"
  },
  "RSI": {
    term: "RSI",
    explanation: "相对强弱指标，衡量价格变动的速度和变化，判断市场买卖力量的对比。",
    simpleExplanation: "股票的‘力气指标’，看看多头和空头谁的力气更大。"
  },
  "换手率": {
    term: "换手率",
    explanation: "在一定时间内市场中股票转手买卖的频率，是反映股票流通性强弱的指标。",
    simpleExplanation: "股票的‘人气值’，换手率高说明大家买卖很活跃，人气旺。"
  },
  "市盈率": {
    term: "市盈率",
    explanation: "股票价格除以每股收益，反映了投资者为获得一元利润所愿意支付的价格。",
    simpleExplanation: "回本年限，假设公司利润不变，你买入后需要多少年能靠利润回本。"
  }
};

export const findTermsInText = (text: string): string[] => {
  const terms = Object.keys(STOCK_GLOSSARY);
  return terms.filter(term => text.includes(term));
};
