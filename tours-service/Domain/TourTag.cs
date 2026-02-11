namespace TourService.Domain;

public class TourTag
{
    public long TourId { get; set; }
    public Tour Tour { get; set; } = null!;

    public long TagId { get; set; }
    public Tag Tag { get; set; } = null!;
}
