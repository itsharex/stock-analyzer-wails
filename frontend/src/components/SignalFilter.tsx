import React from 'react';
import { Filter, X } from 'lucide-react';

interface SignalFilterProps {
  minScore: number;
  onMinScoreChange: (score: number) => void;
  selectedGrades: string[];
  onGradesChange: (grades: string[]) => void;
  onReset: () => void;
}

const SignalFilter: React.FC<SignalFilterProps> = ({
  minScore,
  onMinScoreChange,
  selectedGrades,
  onGradesChange,
  onReset,
}) => {
  const grades = [
    { value: 'S', label: 'S级', min: 90, color: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/50' },
    { value: 'A', label: 'A级', min: 80, color: 'bg-green-500/20 text-green-400 border-green-500/50' },
    { value: 'B', label: 'B级', min: 70, color: 'bg-blue-500/20 text-blue-400 border-blue-500/50' },
    { value: 'C', label: 'C级', min: 0, color: 'bg-gray-500/20 text-gray-400 border-gray-500/50' },
  ];

  const toggleGrade = (grade: string) => {
    if (selectedGrades.includes(grade)) {
      onGradesChange(selectedGrades.filter(g => g !== grade));
    } else {
      onGradesChange([...selectedGrades, grade]);
    }
  };

  const hasActiveFilters = minScore > 0 || selectedGrades.length > 0;

  return (
    <div className="p-3 border-b border-white/10 bg-[#161b22]/50 space-y-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-xs text-gray-400">
          <Filter className="w-3.5 h-3.5" />
          <span>信号筛选</span>
        </div>
        {hasActiveFilters && (
          <button
            onClick={onReset}
            className="flex items-center gap-1 text-xs text-gray-500 hover:text-gray-300 transition-colors"
          >
            <X className="w-3 h-3" />
            重置
          </button>
        )}
      </div>

      {/* 质量等级筛选 */}
      <div className="space-y-2">
        <label className="text-xs text-gray-500">质量等级</label>
        <div className="flex flex-wrap gap-2">
          {grades.map(grade => (
            <button
              key={grade.value}
              onClick={() => toggleGrade(grade.value)}
              className={`px-2 py-1 rounded text-xs border transition-all ${
                selectedGrades.includes(grade.value) || selectedGrades.length === 0
                  ? grade.color
                  : 'bg-gray-800/50 text-gray-600 border-gray-700'
              }`}
            >
              {grade.label}
              {selectedGrades.includes(grade.value) || selectedGrades.length === 0 ? (
                <span className="ml-1">✓</span>
              ) : null}
            </button>
          ))}
        </div>
      </div>

      {/* 最低评分滑块 */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-xs">
          <label className="text-gray-500">最低 AI 评分</label>
          <span className="text-white font-mono">{minScore}</span>
        </div>
        <input
          type="range"
          min="0"
          max="100"
          step="5"
          value={minScore}
          onChange={(e) => onMinScoreChange(Number(e.target.value))}
          className="w-full h-1.5 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-blue-500"
        />
        <div className="flex justify-between text-[10px] text-gray-600 font-mono">
          <span>0</span>
          <span>50</span>
          <span>100</span>
        </div>
      </div>
    </div>
  );
};

export default SignalFilter;
