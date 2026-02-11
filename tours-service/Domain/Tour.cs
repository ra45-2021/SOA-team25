using System.ComponentModel.DataAnnotations;

namespace TourService.Domain;

public class Tour
{
    public long Id { get; set; }

    [Required, MaxLength(120)]
    public string Name { get; set; } = string.Empty;

    [Required, MaxLength(4000)]
    public string Description { get; set; } = string.Empty;

    public TourDifficulty Difficulty { get; set; }

    public TourStatus Status { get; set; }

    public decimal Price { get; set; }

    [Required, MaxLength(64)]
    public string AuthorId { get; set; } = string.Empty;

    public DateTime CreatedAtUtc { get; set; }
    public DateTime UpdatedAtUtc { get; set; }

    public ICollection<TourTag> TourTags { get; set; } = new List<TourTag>();
}
