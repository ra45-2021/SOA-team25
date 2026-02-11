using TourService.Domain;

namespace TourService.Contracts;

public class CreateTourRequest
{
    public string Name { get; set; } = string.Empty;
    public string Description { get; set; } = string.Empty;
    public TourDifficulty Difficulty { get; set; }
    public List<string> Tags { get; set; } = new();
}
