# Implementation Plan: Decision Pioneer Page Logic & UI

I will implement the "Decision Pioneer" page with full interactivity, real-time signal streaming, and responsive design as requested.

## 1. Backend Implementation (Go)
*   **New API Method**: Add `GetLatestSignals(limit int)` to `App` struct (and `StrategyService`).
    *   Queries `stock_strategy_signals` table for the most recent signals.
    *   Returns enriched data including AI scores and reasons.
*   **Real-time Event**: Update `RunStrategyScan` in `app.go` to emit a Wails event `new_signal` whenever a new signal is verified by AI.

## 2. Frontend Architecture
*   **Type Definitions**: Add `Signal` and `AIAnalysis` interfaces to `frontend/src/types.ts`.
*   **State Management**:
    *   `DecisionPioneerPage` will hold the `selectedSignal` state.
    *   **Left Panel** updates this state on click.
    *   **Middle & Right Panels** react to state changes to fetch/display data.

## 3. Component Implementation
### A. Left Panel: Signal List (`SignalList.tsx`)
*   **Data**: Fetches initial list via `GetLatestSignals` and listens to `new_signal` event.
*   **UI**:
    *   Glassmorphism cards (`bg-gray-900/50 backdrop-blur-md`).
    *   Dynamic AI Score Ring (Green > 80, Orange < 60).
    *   Tags for signal types (e.g., "AI Quant Buy").

### B. Middle Panel: Enhanced K-Line (`EnhancedKLineChart.tsx`)
*   **Optimization**: Wrapped in `React.memo` to prevent unnecessary re-renders.
*   **Markers**: Implements logic to display a breathing "B" icon on the signal date.
*   **Sub-chart**: Displays Net Money Flow (Red/Green bars).

### C. Right Panel: AI Analysis (`AIAnalysisPanel.tsx`)
*   **Animation**: CSS/JS animation for the AI Score dashboard (0 to target score).
*   **Content**:
    *   Chips for `ai_reason` keywords.
    *   Progress bars for "Main Concentration" and "Retail Profit".
    *   "Add to Watchlist" button calling backend API.

## 4. Execution Steps
1.  **Backend**: Implement `GetLatestSignals` and event emission.
2.  **Frontend Types**: Define necessary interfaces.
3.  **Components**: Build `SignalList`, `AIAnalysisPanel`, and integrate into `DecisionPioneerPage`.
4.  **Styling**: Apply Figma-matched styles and Tailwind glassmorphism.
5.  **Mocking**: Use mock data initially for UI verification as requested.
