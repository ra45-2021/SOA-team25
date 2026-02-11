using TourService.Domain;

namespace TourService.Contracts;

public class TourDto
{
    public long Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public TourDifficulty Difficulty { get; set; }
    public TourStatus Status { get; set; }
    public decimal Price { get; set; }
    public List<string> Tags { get; set; } = new();
    public DateTime CreatedAtUtc { get; set; }
    public DateTime UpdatedAtUtc { get; set; }
}
